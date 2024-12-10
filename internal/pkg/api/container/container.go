package container

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/containers/image/v5/types"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/container"
	"github.com/warewulf/warewulf/internal/pkg/kernel"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func ContainerBuild(cbp *wwapiv1.ContainerBuildParameter) (err error) {
	if cbp == nil {
		return fmt.Errorf("ContainerBuildParameter is nil")
	}

	var containers []string

	if cbp.All {
		containers, err = container.ListSources()
	} else {
		containers = cbp.ContainerNames
	}

	if len(containers) == 0 {
		return
	}

	for _, c := range containers {
		if !container.ValidSource(c) {
			return fmt.Errorf("VNFS name does not exist: %s", c)
		}

		err = container.Build(c, cbp.Force)
		if err != nil {
			return fmt.Errorf("could not build container %s: %s", c, err)
		}
	}

	if cbp.Default {
		if len(containers) != 1 {
			return fmt.Errorf("can only set default for one container")
		} else {
			var nodeDB node.NodesYaml
			nodeDB, err = node.New()
			if err != nil {
				return fmt.Errorf("could not open node configuration: %s", err)
			}

			// TODO: Don't loop through profiles, instead have a nodeDB function that goes directly to the map
			profiles, _ := nodeDB.FindAllProfiles()
			for _, profile := range profiles {
				wwlog.Debug("Looking for profile default: %s", profile.Id())
				if profile.Id() == "default" {
					wwlog.Debug("Found profile default, setting container name to: %s", containers[0])
					profile.ContainerName = containers[0]
				}
			}
			// TODO: Need a wrapper and flock around this. Sometimes we restart warewulfd and sometimes we don't.
			err = nodeDB.Persist()
			if err != nil {
				return fmt.Errorf("failed to persist nodedb: %w", err)
			}
			fmt.Printf("Set default profile to container: %s\n", containers[0])
		}
	}
	return
}

func ContainerDelete(cdp *wwapiv1.ContainerDeleteParameter) (err error) {
	if cdp == nil {
		return fmt.Errorf("ContainerDeleteParameter is nil")
	}

	nodeDB, err := node.New()
	if err != nil {
		return fmt.Errorf("could not open nodeDB: %s", err)
	}

	nodes, err := nodeDB.FindAllNodes()
	if err != nil {
		return
	}

ARG_LOOP:
	for i := 0; i < len(cdp.ContainerNames); i++ {
		//_, arg := range args {
		containerName := cdp.ContainerNames[i]
		for _, n := range nodes {
			if n.ContainerName == containerName {
				wwlog.Error("Container is configured for nodes, skipping: %s", containerName)
				continue ARG_LOOP
			}
		}

		if !container.ValidSource(containerName) {
			wwlog.Error("Container name is not a valid source: %s", containerName)
			continue
		}
		err := container.DeleteSource(containerName)
		if err != nil {
			wwlog.Error("Could not remove source: %s", containerName)
		}
		err = container.DeleteImage(containerName)
		if err != nil {
			wwlog.Error("Could not remove image files %s", containerName)
		}

		fmt.Printf("Container has been deleted: %s\n", containerName)
	}

	return
}

func ContainerImport(cip *wwapiv1.ContainerImportParameter) (containerName string, err error) {
	if cip == nil {
		err = fmt.Errorf("NodeAddParameter is nil")
		return
	}

	if cip.Name == "" {
		name := path.Base(cip.Source)
		wwlog.Info("Setting VNFS name: %s", name)
		cip.Name = name
	}
	if !container.ValidName(cip.Name) {
		err = fmt.Errorf("VNFS name contains illegal characters: %s", cip.Name)
		return
	}

	containerName = cip.Name
	fullPath := container.SourceDir(cip.Name)

	// container already exists and should be removed first
	if util.IsDir(fullPath) && cip.Force {
		wwlog.Info("Overwriting existing VNFS")
		err = os.RemoveAll(fullPath)
		if err != nil {
			return
		}
	}

	if util.IsDir(fullPath) {
		if !cip.Update {
			err = fmt.Errorf("VNFS Name exists, specify --force, --update, or choose a different name: %s", cip.Name)
			return
		}
		wwlog.Info("Updating existing VNFS")
	} else if strings.HasPrefix(cip.Source, "docker://") || strings.HasPrefix(cip.Source, "docker-daemon://") ||
		strings.HasPrefix(cip.Source, "file://") || util.IsFile(cip.Source) {
		var sCtx *types.SystemContext
		sCtx, err = getSystemContext(cip.OciNoHttps, cip.OciUsername, cip.OciPassword, cip.Platform)
		if err != nil {
			return
		}

		if util.IsFile(cip.Source) && !filepath.IsAbs(cip.Source) {
			cip.Source, err = filepath.Abs(cip.Source)
			if err != nil {
				err = fmt.Errorf("when resolving absolute path of %s, err: %v", cip.Source, err)
				return
			}
		}
		err = container.ImportDocker(cip.Source, cip.Name, sCtx)
		if err != nil {
			err = fmt.Errorf("could not import image: %s", err.Error())
			_ = container.DeleteSource(cip.Name)
			return
		}
	} else if util.IsDir(cip.Source) {
		err = container.ImportDirectory(cip.Source, cip.Name)
		if err != nil {
			err = fmt.Errorf("could not import image: %s", err.Error())
			_ = container.DeleteSource(cip.Name)
			return
		}
	} else {
		err = fmt.Errorf("invalid dir or uri: %s", cip.Source)
		return
	}

	if cip.SyncUser {
		err = container.SyncUids(cip.Name, true)
		if err != nil {
			err = fmt.Errorf("error in user sync, fix error and run 'syncuser' manually: %s", err)
			return
		}
	}

	if cip.Build {
		wwlog.Info("Building container: %s", cip.Name)
		err = container.Build(cip.Name, true)
		if err != nil {
			err = fmt.Errorf("could not build container %s: %s", cip.Name, err.Error())
			return
		}
	}

	if cip.Default {
		var nodeDB node.NodesYaml
		nodeDB, err = node.New()
		if err != nil {
			err = fmt.Errorf("could not open node configuration: %s", err.Error())
			return
		}

		// TODO: Don't loop through profiles, instead have a nodeDB function that goes directly to the map
		profiles, _ := nodeDB.FindAllProfiles()
		for _, profile := range profiles {
			wwlog.Debug("Looking for profile default: %s", profile.Id())
			if profile.Id() == "default" {
				wwlog.Debug("Found profile default, setting container name to: %s", cip.Name)
				profile.ContainerName = cip.Name
			}
		}
		// TODO: We need this in a function with a flock around it.
		// Also need to understand if the daemon restart is only to
		// reload the config or if there is something more.
		err = nodeDB.Persist()
		if err != nil {
			err = fmt.Errorf("failed to persist nodedb: %w", err)
			return
		}

		wwlog.Info("Set default profile to container: %s", cip.Name)
		err = warewulfd.DaemonReload()
		if err != nil {
			err = fmt.Errorf("failed to reload warewulf daemon: %w", err)
			return
		}
	}
	return
}

func ContainerList() (containerInfo []*wwapiv1.ContainerInfo, err error) {
	var sources []string

	sources, err = container.ListSources()
	if err != nil {
		wwlog.Error("%s", err)
		return
	}

	nodeDB, err := node.New()
	if err != nil {
		wwlog.Error("%s", err)
		return
	}

	nodes, err := nodeDB.FindAllNodes()
	if err != nil {
		wwlog.Error("%s", err)
		return
	}

	nodemap := make(map[string]int)
	for _, n := range nodes {
		nodemap[n.ContainerName]++
	}

	for _, source := range sources {
		if nodemap[source] == 0 {
			nodemap[source] = 0
		}

		wwlog.Debug("Finding kernel version for: %s", source)
		kernel := kernel.FindKernels(source).Default()
		kernelVersion := ""
		if kernel != nil {
			kernelVersion = kernel.Version()
		}
		var creationTime uint64
		sourceStat, err := os.Stat(container.SourceDir(source))
		wwlog.Debug("Checking creation time for: %s,%v", container.SourceDir(source), sourceStat.ModTime())
		if err != nil {
			wwlog.Error("%s\n", err)
		} else {
			creationTime = uint64(sourceStat.ModTime().Unix())
		}
		var modTime uint64
		imageStat, err := os.Stat(container.ImageFile(source))
		if err == nil {
			modTime = uint64(imageStat.ModTime().Unix())
		}
		size, err := util.DirSize(container.SourceDir(source))
		if err != nil {
			wwlog.Error("%s\n", err)
		}
		imgSize := container.ImageSize(source)
		imgCSize := container.CompressedImageSize(source)
		containerInfo = append(containerInfo, &wwapiv1.ContainerInfo{
			Name:          source,
			NodeCount:     uint32(nodemap[source]),
			KernelVersion: kernelVersion,
			CreateDate:    creationTime,
			ModDate:       modTime,
			Size:          uint64(size),
			ImgSize:       uint64(imgSize),
			ImgSizeComp:   uint64(imgCSize),
		})

	}
	return
}

func ContainerShow(csp *wwapiv1.ContainerShowParameter) (response *wwapiv1.ContainerShowResponse, err error) {
	containerName := csp.ContainerName

	if !container.ValidName(containerName) {
		err = fmt.Errorf("%s is not a valid container name", containerName)
		return
	}

	rootFsDir := container.RootFsDir(containerName)
	if !util.IsDir(rootFsDir) {
		err = fmt.Errorf("%s is not a valid container", containerName)
		return
	}
	kernel := kernel.FindKernels(containerName).Default()
	kernelVersion := ""
	if kernel != nil {
		kernelVersion = kernel.Version()
	}

	nodeDB, err := node.New()
	if err != nil {
		return
	}

	nodes, err := nodeDB.FindAllNodes()
	if err != nil {
		return
	}

	var nodeList []string
	for _, n := range nodes {
		if n.ContainerName == containerName {
			nodeList = append(nodeList, n.Id())
		}
	}

	response = &wwapiv1.ContainerShowResponse{
		Name:          containerName,
		Rootfs:        rootFsDir,
		Nodes:         nodeList,
		KernelVersion: kernelVersion,
	}
	return
}

func ContainerRename(crp *wwapiv1.ContainerRenameParameter) (err error) {
	// rename the container source folder
	sourceDir := container.SourceDir(crp.ContainerName)
	destDir := container.SourceDir(crp.TargetName)
	err = os.Rename(sourceDir, destDir)
	if err != nil {
		return err
	}

	err = container.DeleteImage(crp.ContainerName)
	if err != nil {
		wwlog.Warn("Could not remove image files for %s: %s", crp.ContainerName, err)
	}

	if crp.Build {
		err = container.Build(crp.TargetName, true)
		if err != nil {
			return err
		}
	}

	// update the nodes profiles container name
	nodeDB, err := node.New()
	if err != nil {
		return err
	}

	nodes, err := nodeDB.FindAllNodes()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if node.ContainerName == crp.ContainerName {
			node.ContainerName = crp.TargetName
		}
	}

	profiles, err := nodeDB.FindAllProfiles()
	if err != nil {
		return err
	}
	for _, profile := range profiles {
		if profile.ContainerName == crp.ContainerName {
			profile.ContainerName = crp.TargetName
		}
	}

	err = nodeDB.Persist()
	if err != nil {
		return err
	}

	err = warewulfd.DaemonStatus()
	if err != nil {
		// warewulfd is not running, skip
		return nil
	}

	// else reload daemon to apply new changes
	return warewulfd.DaemonReload()
}

// create the system context and reading out environment variables
func getSystemContext(noHttps bool, username string, password string, platform string) (sCtx *types.SystemContext, err error) {
	sCtx = &types.SystemContext{}
	// only check env if noHttps wasn't set
	if !noHttps {
		val, ok := os.LookupEnv("WAREWULF_OCI_NOHTTPS")
		if ok {

			noHttps, err = strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("while parsing insecure http option: %v", err)
			}

		}
		// only set this if we want to disable, otherwise leave as undefined
		if noHttps {
			sCtx.DockerInsecureSkipTLSVerify = types.NewOptionalBool(true)
		}
		sCtx.OCIInsecureSkipTLSVerify = noHttps
	}
	if username == "" {
		username, _ = os.LookupEnv("WAREWULF_OCI_USERNAME")
	}
	if password == "" {
		password, _ = os.LookupEnv("WAREWULF_OCI_PASSWORD")
	}
	if username != "" || password != "" {
		if username != "" && password != "" {
			sCtx.DockerAuthConfig = &types.DockerAuthConfig{
				Username: username,
				Password: password,
			}
		} else {
			return nil, fmt.Errorf("oci username and password env vars must be specified together")
		}
	}
	if platform == "" {
		platform, _ = os.LookupEnv("WAREWULF_OCI_PLATFORM")
	}
	if platform != "" {
		sCtx.ArchitectureChoice = platform
	}
	return sCtx, nil
}
