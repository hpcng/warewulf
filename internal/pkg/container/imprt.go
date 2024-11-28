package container

import (
	"context"
	"os"
	"path"

	"github.com/containers/image/v5/types"
	"github.com/containers/storage/drivers/copy"
	"github.com/containers/storage/pkg/reexec"
	"github.com/pkg/errors"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/oci"
	"github.com/warewulf/warewulf/internal/pkg/util"
)

func ImportDocker(uri string, name string, sCtx *types.SystemContext) error {
	OciBlobCacheDir := path.Join(warewulfconf.Get().Paths.Cachedir + "/warewulf")

	err := os.MkdirAll(OciBlobCacheDir, 0755)
	if err != nil {
		return err
	}

	if !ValidName(name) {
		return errors.New("VNFS name contains illegal characters: " + name)
	}

	fullPath := RootFsDir(name)

	err = os.MkdirAll(fullPath, 0755)
	if err != nil {
		return err
	}

	p, err := oci.NewPuller(
		oci.OptSetBlobCachePath(OciBlobCacheDir),
		oci.OptSetSystemContext(sCtx),
	)
	if err != nil {
		return err
	}

	if _, err := p.GenerateID(context.Background(), uri); err != nil {
		return err
	}

	if err := p.Pull(context.Background(), uri, fullPath); err != nil {
		return err
	}

	return nil
}

func ImportDirectory(uri string, name string) error {
	fullPath := RootFsDir(name)

	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		return err
	}

	if !util.IsDir(uri) {
		return errors.New("Import directory does not exist: " + uri)
	}

	if !util.IsFile(path.Join(uri, "/bin/sh")) {
		return errors.New("Source directory has no /bin/sh: " + uri)
	}
	if reexec.Init() {
		return errors.New("couldn't init reexec")
	}
	err = copy.DirCopy(uri, fullPath, copy.Content, true)
	if err != nil {
		return err
	}

	return nil
}
