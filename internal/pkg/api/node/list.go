package apinode

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/hostlist"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
	"gopkg.in/yaml.v3"
)

/*
NodeList lists all to none of the nodes managed by Warewulf. Returns
a formated string slice, with each line as separate string
*/
func NodeList(nodeGet *wwapiv1.GetNodeList) (nodeList wwapiv1.NodeList, err error) {
	// nil is okay for nodeNames
	nodeDB, err := node.New()
	if err != nil {
		return
	}
	nodes, err := nodeDB.FindAllNodes()
	if err != nil {
		return
	}
	nodeGet.Nodes = hostlist.Expand(nodeGet.Nodes)
	sort.Strings(nodeGet.Nodes)

	if nodeGet.Type == wwapiv1.GetNodeList_Simple {
		nodeList.Output = append(nodeList.Output,
			fmt.Sprintf("%s:=:%s:=:%s", "NODE NAME", "PROFILES", "NETWORK"))
		for _, n := range node.FilterNodeListByName(nodes, nodeGet.Nodes) {
			var netNames []string
			for k := range n.NetDevs {
				netNames = append(netNames, k)
			}
			sort.Strings(netNames)
			nodeList.Output = append(nodeList.Output,
				fmt.Sprintf("%s:=:%s:=:%s", n.Id(), strings.Join(n.Profiles, ","), strings.Join(netNames, ", ")))
		}
	} else if nodeGet.Type == wwapiv1.GetNodeList_Network {
		nodeList.Output = append(nodeList.Output,
			fmt.Sprintf("%s:=:%s:=:%s:=:%s:=:%s:=:%s", "NODE", "NETWORK", "HWADDR", "IPADDR", "GATEWAY", "DEVICE"))
		for _, n := range node.FilterNodeListByName(nodes, nodeGet.Nodes) {
			if len(n.NetDevs) > 0 {
				for name := range n.NetDevs {
					nodeList.Output = append(nodeList.Output,
						fmt.Sprintf("%s:=:%s:=:%s:=:%s:=:%s:=:%s", n.Id(), name,
							n.NetDevs[name].Hwaddr,
							n.NetDevs[name].Ipaddr,
							n.NetDevs[name].Gateway,
							n.NetDevs[name].Device))
				}
			} else {
				nodeList.Output = append(nodeList.Output,
					fmt.Sprintf("%s:=:%s:=:%s:=:%s:=:%s:=:%s", n.Id(), "--", "--", "--", "--", "--"))
			}
		}
	} else if nodeGet.Type == wwapiv1.GetNodeList_Ipmi {
		nodeList.Output = append(nodeList.Output,
			fmt.Sprintf("%s:=:%s:=:%s:=:%s:=:%s", "NODE", "IPMI IPADDR", "IPMI PORT", "IPMI USERNAME", "IPMI INTERFACE"))
		for _, n := range node.FilterNodeListByName(nodes, nodeGet.Nodes) {
			nodeList.Output = append(nodeList.Output,
				fmt.Sprintf("%s:=:%s:=:%s:=:%s:=:%s",
					n.Id(),
					n.Ipmi.Ipaddr.String(),
					n.Ipmi.Port,
					n.Ipmi.UserName,
					n.Ipmi.Interface))
		}
	} else if nodeGet.Type == wwapiv1.GetNodeList_Long {
		nodeList.Output = append(nodeList.Output,
			fmt.Sprintf("%s:=:%s:=:%s:=:%s", "NODE NAME", "KERNEL VERSION", "CONTAINER", "OVERLAYS (S/R)"))
		for _, n := range node.FilterNodeListByName(nodes, nodeGet.Nodes) {
			kernelVersion := ""
			if n.Kernel != nil {
				kernelVersion = n.Kernel.Version
			}
			nodeList.Output = append(nodeList.Output,
				fmt.Sprintf("%s:=:%s:=:%s:=:%s", n.Id(),
					kernelVersion,
					n.ContainerName,
					strings.Join(n.SystemOverlay, ",")+"/"+strings.Join(n.RuntimeOverlay, ",")))
		}
	} else if nodeGet.Type == wwapiv1.GetNodeList_All {
		nodeList.Output = append(nodeList.Output,
			fmt.Sprintf("%s:=:%s:=:%s:=:%s", "NODE", "FIELD", "PROFILE", "VALUE"))
		for _, n := range node.FilterNodeListByName(nodes, nodeGet.Nodes) {
			if _, fields, err := nodeDB.MergeNode(n.Id()); err != nil {
				wwlog.Error("unable to merge node %v: %v", n.Id(), err)
				continue
			} else {
				for _, f := range fields.List(n) {
					nodeList.Output = append(nodeList.Output,
						fmt.Sprintf("%s:=:%s:=:%s:=:%s", n.Id(), f.Field, f.Source, f.Value))
				}
			}
		}
	} else if nodeGet.Type == wwapiv1.GetNodeList_YAML || nodeGet.Type == wwapiv1.GetNodeList_JSON {
		filterNodes := node.FilterNodeListByName(nodes, nodeGet.Nodes)
		var buf []byte
		if nodeGet.Type == wwapiv1.GetNodeList_JSON {
			buf, _ = json.MarshalIndent(filterNodes, "", "  ")
		}
		if nodeGet.Type == wwapiv1.GetNodeList_YAML {
			buf, _ = yaml.Marshal(filterNodes)
		}
		nodeList.Output = append(nodeList.Output, string(buf))

	}
	return
}
