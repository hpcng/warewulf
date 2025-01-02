package warewulfd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"github.com/warewulf/warewulf/internal/pkg/container"
	"github.com/warewulf/warewulf/internal/pkg/kernel"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/overlay"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/version"
)

func apiHandler() *web.Service {
	api := web.NewService(openapi3.NewReflector())

	api.OpenAPISchema().SetTitle("Warewulf v4 API")
	api.OpenAPISchema().SetDescription("This service provides an API to a Warewulf v4 server.")
	api.OpenAPISchema().SetVersion(version.GetVersion())

	api.Get("/api/raw-nodes", apiGetRawNodes())
	api.Get("/api/raw-nodes/{id}", apiGetRawNodeByID())
	api.Put("/api/raw-nodes/{id}", apiPutRawNodeByID())

	api.Get("/api/nodes", apiGetNodes())
	api.Get("/api/nodes/{id}", apiGetNodeByID())

	api.Get("/api/profiles", apiGetProfiles())
	api.Get("/api/profiles/{id}", apiGetProfileByID())

	api.Get("/api/containers", apiGetContainers())
	api.Get("/api/conatiners/{name}", apiGetContainerByName())

	api.Get("/api/overlays", apiGetOverlays())
	api.Get("/api/overlays/{name}", apiGetOverlayByName())
	api.Get("/api/overlays/{name}/files/{path}", apiGetOverlayFile())

	api.Docs("/api/docs", swgui.New)

	return api
}

func apiGetRawNodes() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, _ struct{}, output *map[string]*node.Node) error {
		if registry, err := node.New(); err != nil {
			return err
		} else {
			*output = registry.Nodes
			return nil
		}
	})
	u.SetTitle("Get raw nodes")
	u.SetDescription("Get all nodes, without merging in values from associated profiles.")
	u.SetTags("Node")
	return u
}

func apiGetRawNodeByID() usecase.Interactor {
	type getNodeByIDInput struct {
		ID string `path:"id"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input getNodeByIDInput, output *node.Node) error {
		if registry, err := node.New(); err != nil {
			return err
		} else {
			if node_, err := registry.GetNodeOnly(input.ID); err != nil {
				return status.Wrap(fmt.Errorf("node not found: %v (%v)", input.ID, err), status.NotFound)
			} else {
				*output = node_
				return nil
			}
		}
	})
	u.SetTitle("Get a raw node")
	u.SetDescription("Get a node by its ID, without merging in values from associated profiles.")
	u.SetTags("Node")
	u.SetExpectedErrors(status.NotFound)
	return u
}

func apiPutRawNodeByID() usecase.Interactor {
	type putNodeByIDInput struct {
		ID   string    `path:"id"`
		Node node.Node `json:"node"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input putNodeByIDInput, output *node.Node) error {
		if registry, err := node.New(); err != nil {
			return fmt.Errorf("error accessing node registry: %v", err)
		} else {
			if _, ok := registry.Nodes[input.ID]; !ok {
				_, _ = registry.AddNode(input.ID)
			}
			if err := registry.SetNode(input.ID, input.Node); err != nil {
				return fmt.Errorf("error setting node: %v (%v)", input.ID, err)
			} else {
				if node_, err := registry.GetNodeOnly(input.ID); err != nil {
					return fmt.Errorf("node not found after set: %v (%v)", input.ID, err)
				} else {
					*output = node_
					if err := registry.Persist(); err != nil {
						return fmt.Errorf("error persisting node registry: %v", err)
					}
					return nil
				}
			}
		}
	})
	u.SetTitle("Add or update a raw node")
	u.SetDescription("Add or update a raw node and get the resultant node without merging in values from associated profiles.")
	u.SetTags("Node")

	return u
}

func apiGetNodes() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, _ struct{}, output *map[string]*node.Node) error {
		if registry, err := node.New(); err != nil {
			return err
		} else {
			nodeMap := make(map[string]*node.Node)
			if nodeList, err := registry.FindAllNodes(); err != nil {
				return err
			} else {
				for _, n := range nodeList {
					nodeMap[n.Id()] = &n
				}
				*output = nodeMap
				return nil
			}
		}
	})
	u.SetTitle("Get nodes")
	u.SetDescription("Get all nodes, including values from associated profiles.")
	u.SetTags("Node")
	return u
}

func apiGetNodeByID() usecase.Interactor {
	type getNodeByIDInput struct {
		ID string `path:"id"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input getNodeByIDInput, output *node.Node) error {
		if registry, err := node.New(); err != nil {
			return err
		} else {
			if node_, err := registry.GetNode(input.ID); err != nil {
				return status.Wrap(fmt.Errorf("node not found: %v (%v)", input.ID, err), status.NotFound)
			} else {
				*output = node_
				return nil
			}
		}
	})
	u.SetTitle("Get a node")
	u.SetDescription("Get a node by its ID, including values from associated profiles.")
	u.SetTags("Node")
	u.SetExpectedErrors(status.NotFound)
	return u
}

func apiGetProfiles() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, _ struct{}, output *map[string]*node.Profile) error {
		if registry, err := node.New(); err != nil {
			return err
		} else {
			*output = registry.NodeProfiles
			return nil
		}
	})
	u.SetTitle("Get node profiles")
	u.SetDescription("Get all node profiles.")
	u.SetTags("Profile")
	return u
}

func apiGetProfileByID() usecase.Interactor {
	type getProfileByIDInput struct {
		ID string `path:"id"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input getProfileByIDInput, output *node.Profile) error {
		if registry, err := node.New(); err != nil {
			return err
		} else {
			if profile, err := registry.GetProfile(input.ID); err != nil {
				return status.Wrap(fmt.Errorf("profile not found: %v (%v)", input.ID, err), status.NotFound)
			} else {
				*output = profile
				return nil
			}
		}
	})
	u.SetTitle("Get a node profile")
	u.SetDescription("Get a node profile by its ID.")
	u.SetTags("Profile")
	u.SetExpectedErrors(status.NotFound)
	return u
}

type Container struct {
	Kernels  []string `json:"kernels"`
	Size     int      `json:"size"`
	Writable bool     `json:"writable"`
}

func NewContainer(name string) *Container {
	c := new(Container)
	c.Kernels = []string{}
	for _, k := range kernel.FindKernels(name) {
		c.Kernels = append(c.Kernels, k.Path)
	}
	c.Size = container.ImageSize(name)
	c.Writable = container.IsWriteAble(name)
	return c
}

func apiGetContainers() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, _ struct{}, output *map[string]*Container) error {
		m := make(map[string]*Container)
		if names, err := container.ListSources(); err != nil {
			return err
		} else {
			for _, name := range names {
				m[name] = NewContainer(name)
			}
			*output = m
			return nil
		}
	})
	u.SetTitle("Get container images")
	u.SetDescription("Get all container images.")
	u.SetTags("Container")
	return u
}

func apiGetContainerByName() usecase.Interactor {
	type getContainerByNameInput struct {
		Name string `path:"name"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input getContainerByNameInput, output *Container) error {
		if !container.ValidSource(input.Name) {
			return fmt.Errorf("container not found: %v", input.Name)
		} else {
			*output = *NewContainer(input.Name)
			return nil
		}
	})
	u.SetTitle("Get a container")
	u.SetDescription("Get a container by its name.")
	u.SetTags("Container")
	return u
}

type Overlay struct {
	Files []string `json:"files"`
}

func NewOverlay(name string) *Overlay {
	o := new(Overlay)
	o.Files = []string{}
	if files, err := overlay.OverlayGetFiles(name); err == nil {
		o.Files = files
	}
	return o
}

func apiGetOverlays() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, _ struct{}, output *map[string]*Overlay) error {
		m := make(map[string]*Overlay)
		if names, err := overlay.FindOverlays(); err != nil {
			return err
		} else {
			for _, name := range names {
				m[name] = NewOverlay(name)
			}
			*output = m
			return nil
		}
	})
	u.SetTitle("Get overlays")
	u.SetDescription("Get all overlays.")
	u.SetTags("Overlay")
	return u
}

func apiGetOverlayByName() usecase.Interactor {
	type getOverlayByNameInput struct {
		Name string `path:"name"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input getOverlayByNameInput, output *Overlay) error {
		if !overlay.GetOverlay(input.Name).Exists() {
			return fmt.Errorf("overlay not found: %v", input.Name)
		} else {
			*output = *NewOverlay(input.Name)
			return nil
		}
	})
	u.SetTitle("Get an overlay")
	u.SetDescription("Get an overlay by its name.")
	u.SetTags("Overlay")
	return u
}

type OverlayFile struct {
	Overlay  string `json:"overlay"`
	Path     string `json:"path"`
	Contents string `json:"contents"`
}

func (this *OverlayFile) FullPath() string {
	return path.Join(overlay.GetOverlay(this.Overlay).Rootfs(), this.Path)
}

func (this *OverlayFile) Exists() bool {
	return overlay.GetOverlay(this.Overlay).Exists() && util.IsFile(this.FullPath())
}

func (this *OverlayFile) readContents() (string, error) {
	f, err := os.ReadFile(this.FullPath())
	return string(f), err
}

func NewOverlayFile(name string, path string) (*OverlayFile, error) {
	f := new(OverlayFile)
	f.Overlay = name
	f.Path = path
	if contents, err := f.readContents(); err != nil {
		return f, err
	} else {
		f.Contents = contents
		return f, nil
	}
}

func apiGetOverlayFile() usecase.Interactor {
	type getOverlayByNameInput struct {
		Name string `path:"name"`
		Path string `path:"path"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input getOverlayByNameInput, output *OverlayFile) error {
		if relPath, err := url.QueryUnescape(input.Path); err != nil {
			return fmt.Errorf("failed to decode path: %v (%v)", input.Path, err)
		} else {
			if overlayFile, err := NewOverlayFile(input.Name, relPath); err != nil {
				return fmt.Errorf("unable to read overlay file %v:%v", input.Name, relPath)
			} else {
				*output = *overlayFile
				return nil
			}
		}
	})
	u.SetTitle("Get an overlay file")
	u.SetDescription("Get an overlay file by its name and path.")
	u.SetTags("Overlay")
	return u
}
