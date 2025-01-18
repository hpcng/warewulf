package overlay

import (
	"bytes"
	"encoding/gob"
	"os"
	"strconv"
	"time"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/kernel"
	"github.com/warewulf/warewulf/internal/pkg/node"
)

/*
struct which contains the variables to which are available in
the templates.
*/
type TemplateStruct struct {
	Id            string
	Hostname      string
	BuildHost     string
	BuildTime     string
	BuildTimeUnix string
	BuildSource   string
	Ipaddr        string
	Ipaddr6       string
	Netmask       string
	Network       string
	NetworkCIDR   string
	Overlay       string
	Ipv6          bool
	Dhcp          warewulfconf.DHCPConf
	Nfs           warewulfconf.NFSConf
	Ssh           warewulfconf.SSHConf
	Warewulf      warewulfconf.WarewulfConf
	Tftp          warewulfconf.TFTPConf
	Paths         warewulfconf.BuildConfig
	AllNodes      []node.Node
	node.Node
	// backward compatiblity
	Container string
	ThisNode  *node.Node
}

/*
Initialize an TemplateStruct with the given node.NodeInfo
*/
func InitStruct(overlayName string, nodeData node.Node, allNodes []node.Node) (TemplateStruct, error) {
	var tstruct TemplateStruct
	tstruct.Overlay = overlayName
	hostname, _ := os.Hostname()
	tstruct.BuildHost = hostname
	controller := warewulfconf.Get()
	tstruct.ThisNode = &nodeData
	if tstruct.ThisNode.Kernel == nil {
		tstruct.ThisNode.Kernel = new(node.KernelConf)
	}
	if tstruct.ThisNode.Kernel.Version == "" {
		if kernel_ := kernel.FromNode(tstruct.ThisNode); kernel_ != nil {
			tstruct.ThisNode.Kernel.Version = kernel_.Version()
		}
	}
	tstruct.Nfs = *controller.NFS
	tstruct.Ssh = *controller.SSH
	tstruct.Dhcp = *controller.DHCP
	tstruct.Tftp = *controller.TFTP
	tstruct.Paths = *controller.Paths
	tstruct.Warewulf = *controller.Warewulf
	tstruct.Ipaddr = controller.Ipaddr
	tstruct.Ipaddr6 = controller.Ipaddr6
	tstruct.Netmask = controller.Netmask
	tstruct.Network = controller.Network
	// init some convenience vars
	tstruct.Id = nodeData.Id()
	tstruct.Hostname = nodeData.Id()
	tstruct.Container = nodeData.ContainerName
	// Backwards compatibility for templates using "Keys"
	tstruct.AllNodes = allNodes
	dt := time.Now()
	tstruct.BuildTime = dt.Format("01-02-2006 15:04:05 MST")
	tstruct.BuildTimeUnix = strconv.FormatInt(dt.Unix(), 10)
	// tstruct.Node.Tags = map[string]string{}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	if err := enc.Encode(nodeData); err != nil {
		return tstruct, err
	}
	if err := dec.Decode(&tstruct); err != nil {
		return tstruct, err
	}
	return tstruct, nil
}
