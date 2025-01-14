package overlay

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"math/rand"

	"github.com/Masterminds/sprig/v3"

	"github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

var (
	ErrDoesNotExist = fmt.Errorf("overlay does not exist")
)

// Overlay represents an overlay directory path.
type Overlay string

// Name returns the base name of the overlay directory.
//
// This is derived from the full path of the overlay.
func (this Overlay) Name() string {
	return path.Base(this.Path())
}

// Path returns the string representation of the overlay path.
//
// This method allows the Overlay type to be easily converted back to its
// underlying string representation.
func (this Overlay) Path() string {
	return string(this)
}

// Rootfs returns the path to the root filesystem (rootfs) within the overlay.
//
// If the "rootfs" directory exists inside the overlay path, it returns the
// path to the "rootfs" directory. Otherwise, it checks if the overlay path
// itself is a directory and returns that. If neither exists, it defaults to
// returning the "rootfs" path.
func (this Overlay) Rootfs() string {
	rootfs := path.Join(this.Path(), "rootfs")
	if util.IsDir(rootfs) {
		return rootfs
	} else if util.IsDir(this.Path()) {
		return this.Path()
	} else {
		return rootfs
	}
}

// File constructs a full path to a file within the overlay's root filesystem.
//
// Parameters:
//   - filePath: The relative path of the file within the overlay.
//
// Returns:
//   - The full path to the specified file in the overlay's rootfs.
func (this Overlay) File(filePath string) string {
	return path.Join(this.Rootfs(), filePath)
}

// Exists checks whether the overlay path exists and is a directory.
//
// Returns:
//   - true if the overlay path exists and is a directory; false otherwise.
func (this Overlay) Exists() bool {
	return util.IsDir(this.Path())
}

// IsSiteOverlay determines whether the overlay is a site overlay.
//
// A site overlay is identified by its parent directory matching the configured
// site overlay directory path.
//
// Returns:
//   - true if the overlay is a site overlay; false otherwise.
func (this Overlay) IsSiteOverlay() bool {
	return path.Dir(this.Path()) == config.Get().Paths.SiteOverlaydir()
}

// IsDistributionOverlay determines whether the overlay is a distribution overlay.
//
// A distribution overlay is identified by its parent directory matching the configured
// distribution overlay directory path.
//
// Returns:
//   - true if the overlay is a distribution overlay; false otherwise.
func (this Overlay) IsDistributionOverlay() bool {
	return path.Dir(this.Path()) == config.Get().Paths.DistributionOverlaydir()
}

func BuildAllOverlays(nodes []node.Node, allNodes []node.Node, workerCount int) error {
	nodeChan := make(chan node.Node, len(nodes))
	errChan := make(chan error, len(nodes)*2)

	var wg sync.WaitGroup
	worker := func() {
		for n := range nodeChan {
			wwlog.Info("Building system overlays for %s: [%s]", n.Id(), strings.Join(n.SystemOverlay, ", "))
			if err := BuildOverlay(n, allNodes, "system", n.SystemOverlay); err != nil {
				errChan <- fmt.Errorf("could not build system overlays %v for node %s: %w", n.SystemOverlay, n.Id(), err)
			}

			wwlog.Info("Building runtime overlays for %s: [%s]", n.Id(), strings.Join(n.RuntimeOverlay, ", "))
			if err := BuildOverlay(n, allNodes, "runtime", n.RuntimeOverlay); err != nil {
				errChan <- fmt.Errorf("could not build runtime overlays %v for node %s: %w", n.RuntimeOverlay, n.Id(), err)
			}
		}
		wg.Done()
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker()
	}
	for _, n := range nodes {
		nodeChan <- n
	}
	close(nodeChan)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}
	return nil
}

func BuildSpecificOverlays(nodes []node.Node, allNodes []node.Node, overlayNames []string, workerCount int) error {
	nodeChan := make(chan node.Node, len(nodes))
	errChan := make(chan error, len(nodes))

	var wg sync.WaitGroup
	worker := func() {
		for n := range nodeChan {
			wwlog.Info("Building overlay for %s: %v", n, overlayNames)
			for _, overlayName := range overlayNames {
				err := BuildOverlay(n, allNodes, "", []string{overlayName})
				if err != nil {
					errChan <- fmt.Errorf("could not build overlay %s for node %s: %w", overlayName, n.Id(), err)
				}
			}
		}
		wg.Done()
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker()
	}
	for _, n := range nodes {
		nodeChan <- n
	}
	close(nodeChan)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}
	return nil
}

/*
Build overlay for the host, so no argument needs to be given
*/
func BuildHostOverlay() error {
	hostname, _ := os.Hostname()
	hostData := node.NewNode(hostname)
	wwlog.Info("Building overlay for %s: host", hostname)
	hostdir := GetOverlay("host").Rootfs()
	stats, err := os.Stat(hostdir)
	if err != nil {
		return fmt.Errorf("could not build host overlay: %w ", err)
	}
	if !(stats.Mode() == os.FileMode(0750|os.ModeDir) || stats.Mode() == os.FileMode(0700|os.ModeDir)) {
		wwlog.SecWarn("Permissions of host overlay dir %s are %s (750 is considered as secure)", hostdir, stats.Mode())
	}
	registry, err := node.New()
	if err != nil {
		return err
	}
	var allNodes []node.Node
	allNodes, err = registry.FindAllNodes()
	if err != nil {
		return err
	}
	return BuildOverlayIndir(hostData, allNodes, []string{"host"}, "/")
}

/*
Get all overlays present in warewulf
*/
func FindOverlays() (overlayList []string, err error) {
	dotfilecheck, _ := regexp.Compile(`^\..*`)
	controller := config.Get()
	var files []fs.DirEntry
	if distfiles, err := os.ReadDir(controller.Paths.DistributionOverlaydir()); err == nil {
		files = append(files, distfiles...)
	}
	if sitefiles, err := os.ReadDir(path.Join(controller.Paths.SiteOverlaydir())); err == nil {
		files = append(files, sitefiles...)
	}
	for _, file := range files {
		wwlog.Debug("Evaluating overlay source: %s", file.Name())
		isdotfile := dotfilecheck.MatchString(file.Name())

		if (file.IsDir()) && !(isdotfile) {
			overlayList = append(overlayList, file.Name())
		}
	}
	return overlayList, nil
}

/*
Build the given overlays for a node and create a Image for them
*/
func BuildOverlay(nodeConf node.Node, allNodes []node.Node, context string, overlayNames []string) error {
	if len(overlayNames) == 0 && context == "" {
		return nil
	}

	// create the dir where the overlay images will reside
	name := fmt.Sprintf("overlay %s/%v", nodeConf.Id(), overlayNames)
	overlayImage := OverlayImage(nodeConf.Id(), context, overlayNames)
	overlayImageDir := path.Dir(overlayImage)

	err := os.MkdirAll(overlayImageDir, 0750)
	if err != nil {
		return fmt.Errorf("failed to create directory for %s: %s: %w", name, overlayImageDir, err)
	}

	wwlog.Debug("Created directory for %s: %s", name, overlayImageDir)

	buildDir, err := os.MkdirTemp(os.TempDir(), ".wwctl-overlay-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory for %s: %w", name, err)
	}
	defer os.RemoveAll(buildDir)

	wwlog.Debug("Created temporary directory for %s: %s", name, buildDir)

	err = BuildOverlayIndir(nodeConf, allNodes, overlayNames, buildDir)
	if err != nil {
		return fmt.Errorf("failed to generate files for %s: %w", name, err)
	}

	wwlog.Debug("Generated files for %s", name)

	err = util.BuildFsImage(
		name,
		buildDir,
		overlayImage,
		[]string{"*"},
		[]string{},
		// ignore cross-device files
		true,
		"newc")

	return err
}

var regFile *regexp.Regexp
var regLink *regexp.Regexp

func init() {
	regFile = regexp.MustCompile(`.*{{\s*/\*\s*file\s*["'](.*)["']\s*\*/\s*}}.*`)
	regLink = regexp.MustCompile(`.*{{\s*/\*\s*softlink\s*["'](.*)["']\s*\*/\s*}}.*`)
}

// Build the given overlays for a node in the given directory.
func BuildOverlayIndir(nodeData node.Node, allNodes []node.Node, overlayNames []string, outputDir string) error {
	if len(overlayNames) == 0 {
		return nil
	}
	if !util.IsDir(outputDir) {
		return fmt.Errorf("output must a be a directory: %s", outputDir)
	}

	if !util.ValidString(strings.Join(overlayNames, ""), "^[a-zA-Z0-9-._:]+$") {
		return fmt.Errorf("overlay names contains illegal characters: %v", overlayNames)
	}

	wwlog.Verbose("Processing node/overlays: %s/%s", nodeData.Id(), strings.Join(overlayNames, ","))
	for _, overlayName := range overlayNames {
		wwlog.Verbose("Building overlay %s for node %s in %s", overlayName, nodeData.Id(), outputDir)
		overlayRootfs := GetOverlay(overlayName).Rootfs()
		if !util.IsDir(overlayRootfs) {
			return fmt.Errorf("overlay %s: %w", overlayName, ErrDoesNotExist)
		}

		wwlog.Debug("Walking the overlay structure: %s", overlayRootfs)
		err := filepath.Walk(overlayRootfs, func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error for %s: %w", walkPath, err)
			}
			wwlog.Debug("Found overlay file: %s", walkPath)

			relPath, relErr := filepath.Rel(overlayRootfs, walkPath)
			if relErr != nil {
				wwlog.Warn("Error computing relative path for %s: %v", walkPath, relErr)
				return relErr
			}
			outputPath := path.Join(outputDir, relPath)

			if info.IsDir() {
				wwlog.Debug("Found directory: %s", walkPath)

				if err = os.MkdirAll(outputPath, info.Mode()); err != nil {
					return fmt.Errorf("could not create directory within overlay: %w", err)
				}
				if err = util.CopyUIDGID(walkPath, outputPath); err != nil {
					return fmt.Errorf("failed setting permissions on overlay directory: %w", err)
				}

				wwlog.Debug("Created directory in overlay: %s", outputPath)

			} else if filepath.Ext(walkPath) == ".ww" {
				originalOutputPath := outputPath
				outputPath := strings.TrimSuffix(outputPath, ".ww")
				tstruct, err := InitStruct(overlayName, nodeData, allNodes)
				if err != nil {
					return fmt.Errorf("failed to initial data for %s: %w", nodeData.Id(), err)
				}
				tstruct.BuildSource = walkPath
				wwlog.Verbose("Evaluating overlay template file: %s", walkPath)

				buffer, backupFile, writeFile, err := RenderTemplateFile(walkPath, tstruct)
				if err != nil {
					return fmt.Errorf("failed to render template %s: %w", walkPath, err)
				}
				if !writeFile {
					return nil
				}
				var fileBuffer bytes.Buffer
				// search for magic file name comment
				fileScanner := bufio.NewScanner(bytes.NewReader(buffer.Bytes()))
				fileScanner.Split(ScanLines)
				foundFileComment := false
				for fileScanner.Scan() {
					line := fileScanner.Text()
					filenameFromTemplate := regFile.FindAllStringSubmatch(line, -1)
					softlinkFromTemplate := regLink.FindAllStringSubmatch(line, -1)
					if len(softlinkFromTemplate) != 0 {
						wwlog.Debug("Creating soft link %s -> %s", outputPath, softlinkFromTemplate[0][1])
						return os.Symlink(softlinkFromTemplate[0][1], outputPath)
					} else if len(filenameFromTemplate) != 0 {
						wwlog.Debug("Writing file %s", filenameFromTemplate[0][1])
						if foundFileComment {
							err = CarefulWriteBuffer(outputPath, fileBuffer, backupFile, info.Mode())
							if err != nil {
								return fmt.Errorf("could not write file from template: %w", err)
							}
							err = util.CopyUIDGID(walkPath, outputPath)
							if err != nil {
								return fmt.Errorf("failed setting permissions on template output file: %w", err)
							}
							fileBuffer.Reset()
						}
						outputPath = path.Join(path.Dir(originalOutputPath), filenameFromTemplate[0][1])
						foundFileComment = true
					} else {
						if _, err = fileBuffer.WriteString(line); err != nil {
							return fmt.Errorf("could not write to template buffer: %w", err)
						}
					}
				}
				err = CarefulWriteBuffer(outputPath, fileBuffer, backupFile, info.Mode())
				if err != nil {
					return fmt.Errorf("could not write file from template: %w", err)
				}
				err = util.CopyUIDGID(walkPath, outputPath)
				if err != nil {
					return fmt.Errorf("failed setting permissions on template output file: %w", err)
				}
				wwlog.Debug("Wrote template file into overlay: %s", outputPath)

			} else if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				wwlog.Debug("Found symlink %s", walkPath)
				target, err := os.Readlink(walkPath)
				if err != nil {
					return fmt.Errorf("failed reading symlink: %w", err)
				}
				if util.IsFile(outputPath) {
					backupPath := outputPath + ".wwbackup"
					if !util.IsFile(backupPath) {
						wwlog.Debug("Output file already exists: moving to backup file")
						if err = os.Rename(outputPath, backupPath); err != nil {
							return fmt.Errorf("failed renaming to backup file: %w", err)
						}
					} else {
						wwlog.Debug("%s exists, keeping the backup file", backupPath)
						if err = os.Remove(outputPath); err != nil {
							return fmt.Errorf("failed removing existing file: %w", err)
						}
					}
				}
				if err = os.Symlink(target, outputPath); err != nil {
					return fmt.Errorf("failed creating symlink: %w", err)
				}
				wwlog.Debug("Created symlink file: %s", outputPath)
			} else {
				if err := util.CopyFile(walkPath, outputPath); err != nil {
					return fmt.Errorf("could not copy file into overlay: %w", err)
				}
				wwlog.Debug("Copied overlay file: %s", outputPath)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to build overlay image directory: %w", err)
		}
	}

	return nil
}

/*
Writes buffer to the destination file. If wwbackup is set a wwbackup will be created.
*/
func CarefulWriteBuffer(destFile string, buffer bytes.Buffer, backupFile bool, perm fs.FileMode) error {
	wwlog.Debug("Trying to careful write file (%d bytes): %s", buffer.Len(), destFile)
	if backupFile {
		if !util.IsFile(destFile+".wwbackup") && util.IsFile(destFile) {
			err := util.CopyFile(destFile, destFile+".wwbackup")
			if err != nil {
				return fmt.Errorf("failed to create backup: %s -> %s.wwbackup %w", destFile, destFile, err)
			}
		}

	}
	w, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("could not open new file for template %w", err)
	}
	defer w.Close()
	_, err = buffer.WriteTo(w)
	return err
}

/*
Parses the template with the given filename, variables must be in data. Returns the
parsed template as bytes.Buffer, and the bool variables for backupFile and writeFile.
If something goes wrong an error is returned.
*/
func RenderTemplateFile(fileName string, data TemplateStruct) (
	buffer bytes.Buffer,
	backupFile bool,
	writeFile bool,
	err error) {
	backupFile = true
	writeFile = true
	// Build our FuncMap
	funcMap := template.FuncMap{
		"Include":      templateFileInclude,
		"IncludeFrom":  templateContainerFileInclude,
		"IncludeBlock": templateFileBlock,
		"ImportLink":   importSoftlink,
		"basename":     path.Base,
		"random":       func(a []string) int { return rand.Intn(len(a))},
		"inc":          func(i int) int { return i + 1 },
		"dec":          func(i int) int { return i - 1 },
		"file":         func(str string) string { return fmt.Sprintf("{{ /* file \"%s\" */ }}", str) },
		"softlink":     softlink,
		"readlink":     filepath.EvalSymlinks,
		"IgnitionJson": func() string {
			str := createIgnitionJson(data.ThisNode)
			if str != "" {
				return str
			}
			writeFile = false
			return ""
		},
		"abort": func() string {
			wwlog.Debug("abort file called in %s", fileName)
			writeFile = false
			return ""
		},
		"nobackup": func() string {
			wwlog.Debug("not backup for %s", fileName)
			backupFile = false
			return ""
		},
	}

	// Merge sprig.FuncMap with our FuncMap
	for key, value := range sprig.TxtFuncMap() {
		funcMap[key] = value
	}

	// Create the template with the merged FuncMap
	tmpl, err := template.New(path.Base(fileName)).Option("missingkey=default").Funcs(funcMap).ParseGlob(fileName)
	if err != nil {
		err = fmt.Errorf("could not parse template %s: %w", fileName, err)
		return
	}

	err = tmpl.Execute(&buffer, data)
	if err != nil {
		err = fmt.Errorf("could not execute template: %w", err)
		return
	}
	return
}

// Simple version of ScanLines, but include the line break
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// Get all the files as a string slice for a given overlay
func OverlayGetFiles(name string) (files []string, err error) {
	baseDir := GetOverlay(name).Rootfs()
	if !util.IsDir(baseDir) {
		err = fmt.Errorf("overlay %s doesn't exist", name)
		return
	}
	err = filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if util.IsFile(path) {
			files = append(files, strings.TrimPrefix(path, baseDir))
		}
		return nil
	})
	return
}
