package node

import (
	"bytes"
	"encoding/gob"
	"os"
	"path"
	"sort"

	"dario.cat/mergo"
	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"

	"gopkg.in/yaml.v3"
)

const NodesConfPath = "warewulf/nodes.conf"

/*
Get the actual path of nodes.conf, with the possibility to prefix
*/
func GetNodesConf(prefix ...string) string {
	var myPath []string
	if len(prefix) == 0 {
		controller := warewulfconf.Get()
		myPath = append(myPath, controller.Paths.Sysconfdir)
	} else {
		myPath = append(myPath, prefix...)
	}
	myPath = append(myPath, NodesConfPath)
	return path.Join(myPath...)
}

/*
Creates a new nodeDb object from the on-disk configuration
*/
func New() (NodesYaml, error) {
	wwlog.Verbose("Opening node configuration file: %s", GetNodesConf())
	data, err := os.ReadFile(GetNodesConf())
	if err != nil {
		return NodesYaml{}, err
	}
	return Parse(data)
}

// Parse constructs a new nodeDb object from an input YAML
// document. Passes any errors return from yaml.Unmarshal. Returns an
// error if any parsed value is not of a valid type for the given
// parameter.
func Parse(data []byte) (nodeList NodesYaml, err error) {
	wwlog.Debug("Unmarshaling the node configuration")
	err = yaml.Unmarshal(data, &nodeList)
	if err != nil {
		return nodeList, err
	}
	wwlog.Debug("Checking nodes for types")
	if nodeList.Nodes == nil {
		nodeList.Nodes = map[string]*Node{}
	}
	if nodeList.NodeProfiles == nil {
		nodeList.NodeProfiles = map[string]*Profile{}
	}
	wwlog.Debug("returning node object")
	return nodeList, nil
}

/*
Get a node with its merged in nodes
*/
func (config *NodesYaml) GetNode(id string) (node Node, err error) {
	if _, ok := config.Nodes[id]; !ok {
		return node, ErrNotFound
	}
	node = EmptyNode()
	// create a deep copy of the node, as otherwise pointers
	// and not their contents is merged
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err = enc.Encode(config.Nodes[id])
	if err != nil {
		return node, err
	}
	err = dec.Decode(&node)
	if err != nil {
		return node, err
	}
	for _, p := range cleanList(config.Nodes[id].Profiles) {
		includedProfile, err := config.GetProfile(p)
		if err != nil {
			wwlog.Warn("profile not found: %s", p)
			continue
		}
		err = mergo.Merge(&node.Profile, includedProfile, mergo.WithAppendSlice)
		if err != nil {
			return node, err
		}
	}
	// finally set no exported values
	node.id = id
	node.valid = true
	if netdev, ok := node.NetDevs[node.PrimaryNetDev]; ok {
		netdev.primary = true
	} else {
		keys := make([]string, 0)
		for k := range node.NetDevs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if len(keys) > 0 {
			wwlog.Debug("%s: no primary defined, sanitizing to: %s", id, keys[0])
			node.NetDevs[keys[0]].primary = true
			node.PrimaryNetDev = keys[0]
		}
	}
	wwlog.Debug("constructed node: %s", id)
	return
}

/*
Return the node with the id string without the merged in nodes, return ErrNotFound
otherwise
*/
func (config *NodesYaml) GetNodeOnly(id string) (node Node, err error) {
	node = EmptyNode()
	if found, ok := config.Nodes[id]; ok {
		return *found, nil
	}
	return node, ErrNotFound
}

/*
Return pointer to the  node with the id string without the merged in nodes, return ErrMotFound
otherwise
*/
func (config *NodesYaml) GetNodeOnlyPtr(id string) (*Node, error) {
	node := EmptyNode()
	if found, ok := config.Nodes[id]; ok {
		return found, nil
	}
	return &node, ErrNotFound
}

/*
Get the profile with id, return ErrNotFound otherwise
*/
func (config *NodesYaml) GetProfile(id string) (profile Profile, err error) {
	if found, ok := config.NodeProfiles[id]; ok {
		found.id = id
		return *found, nil
	}
	return profile, ErrNotFound
}

/*
Get the profile with id, return ErrNotFound otherwise
*/
func (config *NodesYaml) GetProfilePtr(id string) (profile *Profile, err error) {
	if found, ok := config.NodeProfiles[id]; ok {
		found.id = id
		return found, nil
	}
	return profile, ErrNotFound
}

/*
Get the nodes from the loaded configuration. This function also merges
the nodes with the given nodes.
*/
func (config *NodesYaml) FindAllNodes(nodes ...string) (nodeList []Node, err error) {
	if len(nodes) == 0 {
		for n := range config.Nodes {
			nodes = append(nodes, n)
		}
	}
	wwlog.Debug("Finding nodes: %s", nodes)
	for _, nodeId := range nodes {
		node, err := config.GetNode(nodeId)
		if err != nil {
			return nodeList, err
		}
		nodeList = append(nodeList, node)
	}
	sort.Slice(nodeList, func(i, j int) bool {
		if nodeList[i].ClusterName < nodeList[j].ClusterName {
			return true
		} else if nodeList[i].ClusterName == nodeList[j].ClusterName {
			if nodeList[i].id < nodeList[j].id {
				return true
			}
		}
		return false
	})
	return nodeList, nil
}

/*
Return all nodes as ProfileConf
*/
func (config *NodesYaml) FindAllProfiles(nodes ...string) (profileList []Profile, err error) {
	if len(nodes) == 0 {
		for n := range config.NodeProfiles {
			nodes = append(nodes, n)
		}
	}
	wwlog.Debug("Finding nodes: %s", nodes)
	for _, profileId := range nodes {
		node, err := config.GetProfile(profileId)
		if err != nil {
			return profileList, err
		}
		profileList = append(profileList, node)
	}
	sort.Slice(profileList, func(i, j int) bool {
		if profileList[i].ClusterName < profileList[j].ClusterName {
			return true
		} else if profileList[i].ClusterName == profileList[j].ClusterName {
			if profileList[i].id < profileList[j].id {
				return true
			}
		}
		return false
	})

	return profileList, nil
}

/*
Return the names of all available nodes
*/
func (config *NodesYaml) ListAllNodes() []string {
	nodeList := make([]string, len(config.Nodes))
	for name := range config.Nodes {
		nodeList = append(nodeList, name)
	}
	return nodeList
}

/*
Return the names of all available nodes
*/
func (config *NodesYaml) ListAllProfiles() []string {
	var nodeList []string
	for name := range config.NodeProfiles {
		nodeList = append(nodeList, name)
	}
	return nodeList
}

/*
FindDiscoverableNode returns the first discoverable node and an
interface to associate with the discovered interface. If the nodUNDEFe has
a primary interface, it is returned; otherwise, the first interface
without a hardware address is returned.

If no unconfigured node is found, an error is returned.
*/
func (config *NodesYaml) FindDiscoverableNode() (Node, string, error) {

	nodes, _ := config.FindAllNodes()

	for _, node := range nodes {
		if !(node.Discoverable.Bool()) {
			continue
		}
		if _, ok := node.NetDevs[node.PrimaryNetDev]; ok {
			return node, node.PrimaryNetDev, nil
		}
		for netdev, dev := range node.NetDevs {
			if dev.Hwaddr != "" {
				return node, netdev, nil
			}
		}
	}

	return EmptyNode(), "", ErrNoUnconfigured
}
