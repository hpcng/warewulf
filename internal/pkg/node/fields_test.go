package node

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/warewulf/warewulf/internal/pkg/testenv"
)

func Test_getNestedFieldString(t *testing.T) {
	var tests = map[string]struct {
		nodesConf string
		node      string
		field     string
		value     string
	}{
		"comment (simple)": {
			nodesConf: `
nodes:
  n1:
    comment: n1 comment`,
			node:  "n1",
			field: "Comment",
			value: "n1 comment",
		},
		"kernel args (struct)": {
			nodesConf: `
nodes:
  n1:
    kernel:
      args: n1 args`,
			node:  "n1",
			field: "Kernel.Args",
			value: "n1 args",
		},
		"node tag (map)": {
			nodesConf: `
nodes:
  n1:
    tags:
      tag: n1 tag`,
			node:  "n1",
			field: "Tags[tag]",
			value: "n1 tag",
		},
		"system overlay (slice)": {
			nodesConf: `
nodes:
  n1:
    system overlay:
    - no1
    - no2`,
			node:  "n1",
			field: "SystemOverlay",
			value: "no1,no2",
		},
		"netdev tag (map to struct)": {
			nodesConf: `
nodes:
  n1:
    network devices:
      default:
        tags:
          tag: n1 netdev tag`,
			node:  "n1",
			field: "NetDevs[default].Tags[tag]",
			value: "n1 netdev tag",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			env := testenv.New(t)
			defer env.RemoveAll()
			env.WriteFile("/etc/warewulf/nodes.conf", tt.nodesConf)

			registry, regErr := New()
			assert.NoError(t, regErr)
			node := registry.Nodes[tt.node]
			value, err := getNestedFieldString(node, tt.field)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, value)
		})
	}
}

func Test_listFields(t *testing.T) {
	var tests = map[string]struct {
		object interface{}
		fields []string
	}{
		"node": {
			object: Node{
				Profile: Profile{
					Tags: map[string]string{
						"tag": "value",
					},
					NetDevs: map[string]*NetDev{
						"default": &NetDev{
							Tags: map[string]string{
								"nettag": "netvalue",
							},
						},
					},
					Resources: map[string]Resource{
						"resource": "resvalue",
					},
				},
			},
			fields: []string{
				"Discoverable",
				"AssetKey",
				"Profiles",
				"Comment",
				"ClusterName",
				"ContainerName",
				"Ipxe",
				"RuntimeOverlay",
				"SystemOverlay",
				"Kernel.Version",
				"Kernel.Args",
				"Ipmi.UserName",
				"Ipmi.Password",
				"Ipmi.Ipaddr",
				"Ipmi.Gateway",
				"Ipmi.Netmask",
				"Ipmi.Port",
				"Ipmi.Interface",
				"Ipmi.EscapeChar",
				"Ipmi.Write",
				"Ipmi.Template",
				"Init",
				"Root",
				"NetDevs[default].Type",
				"NetDevs[default].OnBoot",
				"NetDevs[default].Device",
				"NetDevs[default].Hwaddr",
				"NetDevs[default].Ipaddr",
				"NetDevs[default].Ipaddr6",
				"NetDevs[default].Prefix",
				"NetDevs[default].Netmask",
				"NetDevs[default].Gateway",
				"NetDevs[default].MTU",
				"NetDevs[default].Tags[nettag]",
				"Tags[tag]",
				"PrimaryNetDev",
				"Resources[resource]",
			},
		},
		"profile": {
			object: Profile{
				Tags: map[string]string{
					"tag": "value",
				},
				NetDevs: map[string]*NetDev{
					"default": &NetDev{
						Tags: map[string]string{
							"nettag": "netvalue",
						},
					},
				},
				Resources: map[string]Resource{
					"resource": "resvalue",
				},
			},
			fields: []string{
				"Profiles",
				"Comment",
				"ClusterName",
				"ContainerName",
				"Ipxe",
				"RuntimeOverlay",
				"SystemOverlay",
				"Kernel.Version",
				"Kernel.Args",
				"Ipmi.UserName",
				"Ipmi.Password",
				"Ipmi.Ipaddr",
				"Ipmi.Gateway",
				"Ipmi.Netmask",
				"Ipmi.Port",
				"Ipmi.Interface",
				"Ipmi.EscapeChar",
				"Ipmi.Write",
				"Ipmi.Template",
				"Init",
				"Root",
				"NetDevs[default].Type",
				"NetDevs[default].OnBoot",
				"NetDevs[default].Device",
				"NetDevs[default].Hwaddr",
				"NetDevs[default].Ipaddr",
				"NetDevs[default].Ipaddr6",
				"NetDevs[default].Prefix",
				"NetDevs[default].Netmask",
				"NetDevs[default].Gateway",
				"NetDevs[default].MTU",
				"NetDevs[default].Tags[nettag]",
				"Tags[tag]",
				"PrimaryNetDev",
				"Resources[resource]",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.fields, listFields(tt.object))
		})
	}
}

func Test_Field(t *testing.T) {
	field := new(Field)
	assert.Equal(t, "", field.Field)
	assert.Equal(t, "", field.Source)
	assert.Equal(t, "", field.Value)

	field.Field = "test"
	assert.Equal(t, "test", field.Field)

	field.Set("", "value1")
	assert.Equal(t, "", field.Source)
	assert.Equal(t, "value1", field.Value)

	field.Set("", "value2")
	assert.Equal(t, "", field.Source)
	assert.Equal(t, "value2", field.Value)

	field.Set("source3", "value3")
	assert.Equal(t, "source3", field.Source)
	assert.Equal(t, "value3", field.Value)

	field.Set("source4", "value4")
	assert.Equal(t, "source4", field.Source)
	assert.Equal(t, "value4", field.Value)

	field.Set("", "value5")
	assert.Equal(t, "SUPERSEDED", field.Source)
	assert.Equal(t, "value5", field.Value)
}

func Test_fieldMap(t *testing.T) {
	fieldMap := make(fieldMap)
	assert.Equal(t, 0, len(fieldMap))

	fieldMap.Set("field", "", "value1")
	assert.Equal(t, "", fieldMap.Source("field"))
	assert.Equal(t, "value1", fieldMap.Value("field"))

	fieldMap.Set("field", "", "value2")
	assert.Equal(t, "", fieldMap.Source("field"))
	assert.Equal(t, "value2", fieldMap.Value("field"))

	fieldMap.Set("field", "source3", "value3")
	assert.Equal(t, "source3", fieldMap.Source("field"))
	assert.Equal(t, "value3", fieldMap.Value("field"))

	fieldMap.Set("field", "source4", "value4")
	assert.Equal(t, "source4", fieldMap.Source("field"))
	assert.Equal(t, "value4", fieldMap.Value("field"))

	fieldMap.Set("field", "", "value5")
	assert.Equal(t, "SUPERSEDED", fieldMap.Source("field"))
	assert.Equal(t, "value5", fieldMap.Value("field"))
}
