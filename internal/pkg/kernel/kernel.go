package kernel

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"

	"github.com/warewulf/warewulf/internal/pkg/container"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

var (
	kernelSearchPaths = []string{
		"/boot/Image-*", // this is the aarch64 for SUSE, vmlinux which is also present won't boot
		"/boot/vmlinuz-linux*",
		"/boot/vmlinuz-*",
		"/lib/modules/*/vmlinuz",
		"/lib/modules/*/vmlinuz.gz",
	}
)

type collection []*Kernel

func (k collection) Len() int {
	return len(k)
}

func (k collection) Less(i, j int) bool {
	iv := k[i].version()
	jv := k[j].version()
	return (iv == nil && jv != nil) ||
		(iv != nil && jv != nil && iv.LessThan(jv))
}

func (k collection) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func (k collection) Sort() {
	sort.Sort(k)
}

func (k collection) Default() *Kernel {
	nk := append(collection{}, k...)
	sort.Sort(sort.Reverse(nk))
	for _, kernel := range nk {
		if !(kernel.IsDebug() || kernel.IsRescue()) {
			return kernel
		}
	}
	return nil
}

func (k collection) Version(version string) *Kernel {
	for _, kernel := range k {
		if kernel.IsDebug() || kernel.IsRescue() {
			continue
		} else if strings.HasPrefix(kernel.Version(), version) {
			return kernel
		}
	}
	return nil
}

type Kernel struct {
	Path          string
	ContainerName string
}

func FromNode(node *node.Node) *Kernel {
	wwlog.Debug("FromNode(%v)", node)
	if node.ContainerName == "" {
		return nil
	} else if node.Kernel != nil && node.Kernel.Version != "" {
		kernel := &Kernel{ContainerName: node.ContainerName, Path: filepath.Join("/", node.Kernel.Version)}
		if util.IsFile(kernel.FullPath()) {
			return kernel
		} else {
			return FindKernels(node.ContainerName).Version(node.Kernel.Version)
		}
	} else {
		return FindKernels(node.ContainerName).Default()
	}
}

func FindKernelsFromPattern(containerName string, pattern string) (kernels collection) {
	wwlog.Debug("FindKernelsFromPattern(%v, %v)", containerName, pattern)
	root := container.RootFsDir(containerName)
	fullPaths, err := filepath.Glob(filepath.Join(root, pattern))
	wwlog.Debug("%v: fullPaths: %v", filepath.Join(root, pattern), fullPaths)
	if err != nil {
		panic(err)
	}
	for _, fullPath := range fullPaths {
		path, err := filepath.Rel(root, fullPath)
		if err != nil {
			continue
		} else {
			kernels = append(kernels, &Kernel{ContainerName: containerName, Path: filepath.Join("/", path)})
		}
	}
	return kernels
}

func FindKernels(containerName string) (kernels collection) {
	wwlog.Debug("FindKernels(%v)", containerName)
	for _, pattern := range kernelSearchPaths {
		kernels = append(kernels, FindKernelsFromPattern(containerName, pattern)...)
	}
	return kernels
}

func FindAllKernels() (kernels collection) {
	wwlog.Debug("FindAllKernels()")
	if sources, err := container.ListSources(); err == nil {
		for _, source := range sources {
			kernels = append(kernels, FindKernels(source)...)
		}
	} else {
		wwlog.Error("%s", err)
	}
	return kernels
}

func (this *Kernel) version() *version.Version {
	return util.ParseVersion(this.Path)
}

func (this *Kernel) Version() string {
	version := this.version()
	if version == nil {
		return ""
	} else {
		return version.String()
	}
}

func (this *Kernel) IsDebug() bool {
	return strings.Contains(this.Path, "+debug")
}

func (this *Kernel) IsRescue() bool {
	return strings.Contains(this.Path, "-rescue")
}

func (this *Kernel) FullPath() string {
	root := container.RootFsDir(this.ContainerName)
	return filepath.Join(root, this.Path)
}
