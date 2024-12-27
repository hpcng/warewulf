package node

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/warewulf/warewulf/internal/pkg/testenv"
)

func Test_getNodeProfiles(t *testing.T) {
	var tests = map[string]struct {
		nodesConf string
		node      string
		profiles  []string
	}{
		"no profiles": {
			nodesConf: `
nodes:
  n1: {}`,
			node:     "n1",
			profiles: nil,
		},

		"one profile": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1`,
			node:     "n1",
			profiles: []string{"p1"},
		},
		"two profiles": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2`,
			node:     "n1",
			profiles: []string{"p1", "p2"},
		},
		"negated profiles": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
nodeprofiles:
  p2:
    profiles:
    - "~p1"`,
			node:     "n1",
			profiles: []string{"p2"},
		},
		"negated missing profile": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
nodeprofiles:
  p2:
    profiles:
    - "~p3"`,
			node:     "n1",
			profiles: []string{"p1", "p2"},
		},
		"single nested profile": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    profiles:
    - p2`,
			node:     "n1",
			profiles: []string{"p1", "p2"},
		},
		"double nested profile": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    profiles:
    - p2
  p2:
    profiles:
    - p3`,
			node:     "n1",
			profiles: []string{"p1", "p2", "p3"},
		},
		"negated nested profile": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    profiles:
    - p2
  p2:
    profiles:
    - "~p2"
    - p3`,
			node:     "n1",
			profiles: []string{"p1", "p3"},
		},
		"cicular nested profile": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    profiles:
    - p2
  p2:
    profiles:
    - p3
  p3:
    profiles:
    - p1`,
			node:     "n1",
			profiles: []string{"p1", "p2", "p3"},
		},
		"cicular nested profile negation": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    profiles:
    - p2
  p2:
    profiles:
    - "~p1"
    - p3
  p3:
    profiles:
    - p1`,
			node:     "n1",
			profiles: []string{"p2", "p3", "p1"},
		},
		"repeated nested profile": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - pa1
    - pb1
nodeprofiles:
  pa1:
    profiles:
    - pa2
  pb1:
    profiles:
    - "~pa2"
    - pb2
  pb2:
    profiles:
    - pa2`,
			node:     "n1",
			profiles: []string{"pa1", "pb1", "pb2", "pa2"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			env := testenv.New(t)
			defer env.RemoveAll()
			env.WriteFile("/etc/warewulf/nodes.conf", tt.nodesConf)

			registry, regErr := New()
			assert.NoError(t, regErr)
			assert.Equal(t, tt.profiles, registry.getNodeProfiles(tt.node))
		})
	}
}

func Test_MergeNode(t *testing.T) {
	var tests = map[string]struct {
		nodesConf string
		node      string
		field     string
		source    string
		value     string
		nodes     []string
		fields    []string
		sources   []string
		values    []string
	}{
		"node comment": {
			nodesConf: `
nodes:
  n1:
    comment: n1 comment`,
			node:   "n1",
			field:  "Comment",
			source: "",
			value:  "n1 comment",
		},
		"profile comment": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    comment: p1 comment`,
			node:   "n1",
			field:  "Comment",
			source: "p1",
			value:  "p1 comment",
		},
		"multiple profile comments": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
nodeprofiles:
  p1:
    comment: p1 comment
  p2:
    comment: p2 comment`,
			node:   "n1",
			field:  "Comment",
			source: "p2",
			value:  "p2 comment",
		},
		"node comment supersedes profile comment": {
			nodesConf: `
nodes:
  n1:
    comment: n1 comment
    profiles:
    - p1
nodeprofiles:
  p1:
    comment: p1 comment`,
			node:   "n1",
			field:  "Comment",
			source: "SUPERSEDED",
			value:  "n1 comment",
		},
		"node comment supersedes multiple profile comments": {
			nodesConf: `
nodes:
  n1:
    comment: n1 comment
    profiles:
    - p1
    - p2
nodeprofiles:
  p1:
    comment: p1 comment
  p2:
    comment: p2 comment`,
			node:   "n1",
			field:  "Comment",
			source: "SUPERSEDED",
			value:  "n1 comment",
		},
		"node kernel args": {
			nodesConf: `
nodes:
  n1:
    kernel:
      args: n1 args`,
			node:   "n1",
			field:  "Kernel.Args",
			source: "",
			value:  "n1 args",
		},
		"profile kernel args": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    kernel:
      args: p1 args`,
			node:   "n1",
			field:  "Kernel.Args",
			source: "p1",
			value:  "p1 args",
		},
		"multiple profile kernel args": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
nodeprofiles:
  p1:
    kernel:
      args: p1 args
  p2:
    kernel:
      args: p2 args`,
			node:   "n1",
			field:  "Kernel.Args",
			source: "p2",
			value:  "p2 args",
		},
		"node kernel args supersedes profile kernel args": {
			nodesConf: `
nodes:
  n1:
    kernel:
      args: n1 args
    profiles:
    - p1
nodeprofiles:
  p1:
    kernel:
      args: p1 args`,
			node:   "n1",
			field:  "Kernel.Args",
			source: "SUPERSEDED",
			value:  "n1 args",
		},
		"node kernel args supersedes multiple profile kernel args": {
			nodesConf: `
nodes:
  n1:
    kernel:
      args: n1 args
    profiles:
    - p1
    - p2
nodeprofiles:
  p1:
    kernel:
      args: p1 args
  p2:
    kernel:
      args: p2 args`,
			node:   "n1",
			field:  "Kernel.Args",
			source: "SUPERSEDED",
			value:  "n1 args",
		},
		"node tag": {
			nodesConf: `
nodes:
  n1:
    tags:
      tag: n1 tag`,
			node:   "n1",
			field:  "Tags[tag]",
			source: "",
			value:  "n1 tag",
		},
		"profile tag": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    tags:
      tag: p1 tag`,
			node:   "n1",
			field:  "Tags[tag]",
			source: "p1",
			value:  "p1 tag",
		},
		"multiple profile tags": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
nodeprofiles:
  p1:
    tags:
      tag: p1 tag
  p2:
    tags:
      tag: p2 tag`,
			node:   "n1",
			field:  "Tags[tag]",
			source: "p2",
			value:  "p2 tag",
		},
		"node tag supersedes profile tag": {
			nodesConf: `
nodes:
  n1:
    tags:
      tag: n1 tag
    profiles:
    - p1
nodeprofiles:
  p1:
    tags:
      tag: p1 tag`,
			node:   "n1",
			field:  "Tags[tag]",
			source: "SUPERSEDED",
			value:  "n1 tag",
		},
		"node tag supersedes multiple profile tags": {
			nodesConf: `
nodes:
  n1:
    tags:
      tag: n1 tag
    profiles:
    - p1
    - p2
nodeprofiles:
  p1:
    tags:
      tag: p1 tag
  p2:
    tags:
      tag: p2 tag`,
			node:   "n1",
			field:  "Tags[tag]",
			source: "SUPERSEDED",
			value:  "n1 tag",
		},
		"mixture of tags from nodes and profiles": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
    tags:
      n1: n1 tag
nodeprofiles:
  p1:
    tags:
      p1: p1 tag
  p2:
    tags:
      p2: p2 tag`,
			nodes: []string{
				"n1",
				"n1",
				"n1",
			},
			fields: []string{
				"Tags[n1]",
				"Tags[p1]",
				"Tags[p2]",
			},
			sources: []string{
				"",
				"p1",
				"p2",
			},
			values: []string{
				"n1 tag",
				"p1 tag",
				"p2 tag",
			},
		},
		"node system overlay": {
			nodesConf: `
nodes:
  n1:
    system overlay:
    - no1
    - no2`,
			node:   "n1",
			field:  "SystemOverlay",
			source: "",
			value:  "no1,no2",
		},
		"profile system overlay": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    system overlay:
    - po1
    - po2`,
			node:   "n1",
			field:  "SystemOverlay",
			source: "p1",
			value:  "po1,po2",
		},
		"multiple profile system overlays": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
nodeprofiles:
  p1:
    system overlay:
    - po1
    - po2
  p2:
    system overlay:
    - po3
    - po4`,
			node:   "n1",
			field:  "SystemOverlay",
			source: "p1,p2",
			value:  "po1,po2,po3,po4",
		},
		"node system overlay adds to profile system overlay": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    system overlay:
    - no1
    - no2
nodeprofiles:
  p1:
    system overlay:
    - po1
    - po2`,
			node:   "n1",
			field:  "SystemOverlay",
			source: "p1,n1",
			value:  "po1,po2,no1,no2",
		},
		"node system overlay adds to multiple profile system overlays": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
    system overlay:
    - no1
    - no2
nodeprofiles:
  p1:
    system overlay:
    - po1
    - po2
  p2:
    system overlay:
    - po3
    - po4`,
			node:   "n1",
			field:  "SystemOverlay",
			source: "p1,p2,n1",
			value:  "po1,po2,po3,po4,no1,no2",
		},
		"node profiles": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2`,
			node:   "n1",
			field:  "Profiles",
			source: "",
			value:  "p1,p2",
		},
		"nested profiles": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    profiles:
    - p2`,
			node:   "n1",
			field:  "Profiles",
			source: "",
			value:  "p1",
		},
		"negated profiles": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
nodeprofiles:
  p2:
    profiles:
    - "~p1"`,
			node:   "n1",
			field:  "Profiles",
			source: "",
			value:  "p1,p2",
		},
		"node netdev tag": {
			nodesConf: `
nodes:
  n1:
    network devices:
      default:
        tags:
          tag: n1 netdev tag`,
			node:   "n1",
			field:  "NetDevs[default].Tags[tag]",
			source: "",
			value:  "n1 netdev tag",
		},
		"profile netdev tag": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
nodeprofiles:
  p1:
    network devices:
      default:
        tags:
          tag: p1 netdev tag`,
			node:   "n1",
			field:  "NetDevs[default].Tags[tag]",
			source: "p1",
			value:  "p1 netdev tag",
		},
		"multiple profile netdev tags": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
nodeprofiles:
  p1:
    network devices:
      default:
        tags:
          tag: p1 netdev tag
  p2:
    network devices:
      default:
        tags:
          tag: p2 netdev tag`,
			node:   "n1",
			field:  "NetDevs[default].Tags[tag]",
			source: "p2",
			value:  "p2 netdev tag",
		},
		"node supercededs profile netdev tag": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    network devices:
      default:
        tags:
          tag: n1 netdev tag
nodeprofiles:
  p1:
    network devices:
      default:
        tags:
          tag: p1 netdev tag`,
			node:   "n1",
			field:  "NetDevs[default].Tags[tag]",
			source: "SUPERSEDED",
			value:  "n1 netdev tag",
		},
		"node supersedes multiple profile netdev tags": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
    network devices:
      default:
        tags:
          tag: n1 netdev tag
nodeprofiles:
  p1:
    network devices:
      default:
        tags:
          tag: p1 netdev tag
  p2:
    network devices:
      default:
        tags:
          tag: p2 netdev tag`,
			node:   "n1",
			field:  "NetDevs[default].Tags[tag]",
			source: "SUPERSEDED",
			value:  "n1 netdev tag",
		},
		"mixture of netdev tags from nodes and profiles": {
			nodesConf: `
nodes:
  n1:
    profiles:
    - p1
    - p2
    network devices:
      default:
        tags:
          n1: n1 netdev tag
nodeprofiles:
  p1:
    network devices:
      default:
        tags:
          p1: p1 netdev tag
  p2:
    network devices:
      default:
        tags:
          p2: p2 netdev tag`,
			nodes: []string{
				"n1",
				"n1",
				"n1",
			},
			fields: []string{
				"NetDevs[default].Tags[n1]",
				"NetDevs[default].Tags[p1]",
				"NetDevs[default].Tags[p2]",
			},
			sources: []string{
				"",
				"p1",
				"p2",
			},
			values: []string{
				"n1 netdev tag",
				"p1 netdev tag",
				"p2 netdev tag",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			env := testenv.New(t)
			defer env.RemoveAll()
			env.WriteFile("/etc/warewulf/nodes.conf", tt.nodesConf)

			registry, regErr := New()
			assert.NoError(t, regErr)

			if tt.node != "" {
				node, fields, mergeErr := registry.MergeNode(tt.node)
				assert.NoError(t, mergeErr)

				value, valueErr := getNestedFieldString(node, tt.field)
				assert.NoError(t, valueErr)
				assert.Equal(t, tt.value, value)
				assert.Equal(t, tt.value, fields.Value(tt.field))
				assert.Equal(t, tt.source, fields.Source(tt.field))
			}

			var nodes []Node
			for i, _ := range tt.nodes {
				node, _, mergeErr := registry.MergeNode(tt.nodes[i])
				assert.NoError(t, mergeErr)
				nodes = append(nodes, node)
			}

			for i, _ := range tt.nodes {
				_, fields, _ := registry.MergeNode(tt.nodes[i])
				value, valueErr := getNestedFieldString(nodes[i], tt.fields[i])
				assert.NoError(t, valueErr)
				assert.Equal(t, tt.values[i], value)
				assert.Equal(t, tt.values[i], fields.Value(tt.fields[i]))
				assert.Equal(t, tt.sources[i], fields.Source(tt.fields[i]))
			}
		})
	}
}

func Test_MergeNodeKernel(t *testing.T) {
	var tests = map[string]struct {
		nodesConf string
		node      string
		kernel    *KernelConf
	}{
		"interference": {
			nodesConf: `
nodeprofiles:
  default: {}
nodes:
  node1: {}
  test:
    profiles:
      - default
    kernel:
      version: v1.0.0
      args: kernel-args`,
			node:   "node1",
			kernel: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			env := testenv.New(t)
			defer env.RemoveAll()
			env.WriteFile("/etc/warewulf/nodes.conf", tt.nodesConf)

			registry, regErr := New()
			assert.NoError(t, regErr)
			node, _, mergeErr := registry.MergeNode(tt.node)
			assert.NoError(t, mergeErr)
			assert.Equal(t, tt.kernel, node.Kernel)
		})
	}
}

func Test_MergeNodeIpmi(t *testing.T) {
	var tests = map[string]struct {
		nodesConf string
		node      string
		ipmi      *IpmiConf
	}{
		"empty single node": {
			nodesConf: `
nodes:
  n1: {}`,
			node: "n1",
			ipmi: nil,
		},
		"empty node among multiple": {
			nodesConf: `
nodes:
  n0: {}
  n1: {}
  n2: {}`,
			node: "n1",
			ipmi: nil,
		},
		"populated node among multiple": {
			nodesConf: `
nodes:
  n0: {}
  n1:
    ipmi:
      username: root
      password: passw0rd
  n2: {}`,
			node: "n1",
			ipmi: &IpmiConf{
				UserName: "root",
				Password: "passw0rd",
			},
		},
		"no previous interference": {
			nodesConf: `
nodes:
  n0: {}
  n1:
    ipmi:
      username: root
      password: passw0rd
  n2: {}`,
			node: "n0",
			ipmi: nil,
		},
		"no later interference": {
			nodesConf: `
nodes:
  n0: {}
  n1:
    ipmi:
      username: root
      password: passw0rd
  n2: {}`,
			node: "n2",
			ipmi: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			env := testenv.New(t)
			defer env.RemoveAll()
			env.WriteFile("/etc/warewulf/nodes.conf", tt.nodesConf)

			registry, regErr := New()
			assert.NoError(t, regErr)
			node, _, mergeErr := registry.MergeNode(tt.node)
			assert.NoError(t, mergeErr)
			assert.Equal(t, tt.ipmi, node.Ipmi)
		})
	}
}
