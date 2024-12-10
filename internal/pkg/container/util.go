package container

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func ValidName(name string) bool {
	if !util.ValidString(name, "^[\\w\\-\\.\\:]+$") {
		wwlog.Warn("VNFS name has illegal characters: %s", name)
		return false
	}
	return true
}

func ListSources() ([]string, error) {
	var ret []string

	err := os.MkdirAll(SourceParentDir(), 0755)
	if err != nil {
		return ret, errors.New("Could not create VNFS source parent directory: " + SourceParentDir())
	}
	wwlog.Debug("Searching for VNFS Rootfs directories: %s", SourceParentDir())

	sources, err := os.ReadDir(SourceParentDir())
	if err != nil {
		return ret, err
	}

	for _, source := range sources {
		wwlog.Verbose("Found VNFS source: %s", source.Name())

		if !ValidName(source.Name()) {
			continue
		}

		if !ValidSource(source.Name()) {
			continue
		}

		ret = append(ret, source.Name())
	}

	return ret, nil
}

func DoesSourceExist(name string) bool {
	fullPath := RootFsDir(name)
	return util.IsDir(fullPath)
}

func ValidSource(name string) bool {
	if !ValidName(name) {
		return false
	}

	if !DoesSourceExist(name) {
		wwlog.Verbose("Location is not a VNFS source directory: %s", name)
		return false
	}

	return true
}

/*
Delete the chroot of a container
*/
func DeleteSource(name string) error {
	fullPath := SourceDir(name)

	wwlog.Verbose("Removing path: %s", fullPath)
	return os.RemoveAll(fullPath)
}

func Duplicate(name string, destination string) error {
	fullPathImageSource := RootFsDir(name)

	wwlog.Info("Copying sources...")
	err := ImportDirectory(fullPathImageSource, destination)

	if err != nil {
		return err
	}

	return nil
}

/*
Delete the image of a container
*/
func DeleteImage(name string) error {
	imageFile := ImageFile(name)
	if util.IsFile(imageFile) {
		wwlog.Verbose("removing %s for container %s", imageFile, name)
		errImg := os.Remove(imageFile)
		wwlog.Verbose("removing %s for container %s", imageFile+".gz", name)
		errGz := os.Remove(imageFile + ".gz")
		if errImg != nil {
			return errors.Errorf("Problems delete %s for container %s: %s\n", imageFile, name, errImg)
		}
		if errGz != nil {
			return errors.Errorf("Problems delete %s for container %s: %s\n", imageFile+".gz", name, errGz)
		}
		return nil
	}
	return errors.Errorf("Image %s of container %s doesn't exist\n", imageFile, name)
}

func IsWriteAble(name string) bool {
	return !util.IsFile(filepath.Join(SourceDir(name), "readonly"))
}

func ImageSize(name string) int {
	if img, err := os.Stat(ImageFile(name)); err == nil {
		return int(img.Size())
	} else {
		return 0
	}
}

func CompressedImageSize(name string) int {
	if img, err := os.Stat(CompressedImageFile(name)); err == nil {
		return int(img.Size())
	} else {
		return 0
	}
}
