package node

import (
	"net"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

type nodeList []Node

func (a nodeList) Len() int           { return len(a) }
func (a nodeList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nodeList) Less(i, j int) bool { return a[i].id < a[j].id }

type profileList []Profile

func (a profileList) Len() int           { return len(a) }
func (a profileList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a profileList) Less(i, j int) bool { return a[i].id < a[j].id }

/**********
 *
 * Filters
 *
 *********/

/*
Filter a given slice of NodeConf against a given
regular expression
*/
func FilterNodeListByName(set []Node, searchList []string) []Node {
	var ret []Node
	unique := make(map[string]Node)

	if len(searchList) > 0 {
		for _, search := range searchList {
			for _, entry := range set {
				if match, _ := regexp.MatchString("^"+search+"$", entry.id); match {
					unique[entry.id] = entry
				}
			}
		}
		for _, n := range unique {
			ret = append(ret, n)
		}
	} else {
		ret = set
	}

	sort.Sort(nodeList(ret))
	return ret
}

/*
Filter a given slice of ProfileConf against a given
regular expression
*/
func FilterProfileListByName(set []Profile, searchList []string) []Profile {
	var ret []Profile
	unique := make(map[string]Profile)

	if len(searchList) > 0 {
		for _, search := range searchList {
			for _, entry := range set {
				if match, _ := regexp.MatchString("^"+search+"$", entry.id); match {
					unique[entry.id] = entry
				}
			}
		}
		for _, n := range unique {
			ret = append(ret, n)
		}
	} else {
		ret = set
	}

	sort.Sort(profileList(ret))
	return ret
}

/*
Creates an NodeConf with the given id. Doesn't add it to the database
*/
func NewNode(id string) (nodeconf Node) {
	nodeconf = EmptyNode()
	nodeconf.id = id
	return nodeconf
}

func NewProfile(id string) (profileconf Profile) {
	profileconf = EmptyProfile()
	profileconf.id = id
	return profileconf
}

func EmptyNode() (nodeconf Node) {
	nodeconf.Ipmi = new(IpmiConf)
	nodeconf.Ipmi.Tags = map[string]string{}
	nodeconf.Kernel = new(KernelConf)
	nodeconf.NetDevs = make(map[string]*NetDev)
	nodeconf.Tags = map[string]string{}
	return nodeconf
}

/*
Creates a ProfileConf but doesn't add it to the database.
*/
func EmptyProfile() (profileconf Profile) {
	profileconf.Ipmi = new(IpmiConf)
	profileconf.Ipmi.Tags = map[string]string{}
	profileconf.Kernel = new(KernelConf)
	profileconf.NetDevs = make(map[string]*NetDev)
	profileconf.Tags = map[string]string{}
	return profileconf
}

/*
Flattens out a NodeConf, which means if there are no explicit values in *IpmiConf
or *KernelConf, these pointer will set to nil. This will remove something like
ipmi: {} from nodes.conf
*/
func (info *Node) Flatten() {
	recursiveFlatten(info)
}

/*
Flattens out a ProfileConf, which means if there are no explicit values in *IpmiConf
or *KernelConf, these pointer will set to nil. This will remove something like
ipmi: {} from nodes.conf
*/
func (info *Profile) Flatten() {
	recursiveFlatten(info)
}

func recursiveFlatten(obj interface{}) (hasContent bool) {
	valObj := reflect.ValueOf(obj)
	typeObj := reflect.TypeOf(obj)
	hasContent = false
	if valObj.IsNil() {
		return
	}
	for i := 0; i < typeObj.Elem().NumField(); i++ {
		if valObj.Elem().Field(i).IsValid() {
			if !typeObj.Elem().Field(i).IsExported() {
				continue
			}
		}
		switch typeObj.Elem().Field(i).Type.Kind() {
		case reflect.Map:
			mapIter := valObj.Elem().Field(i).MapRange()
			for mapIter.Next() {
				switch mapIter.Value().Kind() {
				case reflect.Map, reflect.Pointer, reflect.Slice:
					if mapIter.Value().Type().Elem().Kind() == reflect.Struct {
						ret := recursiveFlatten(mapIter.Value().Interface())
						hasContent = ret || hasContent
					} else {
						hasContent = !mapIter.Value().IsZero() || hasContent
					}
				default:
					hasContent = !mapIter.Value().IsZero() || hasContent
				}
			}

		case reflect.Ptr:
			if valObj.Elem().Field(i).Addr().IsValid() {
				ret := recursiveFlatten((valObj.Elem().Field(i).Interface()))
				if !ret {
					valObj.Elem().Field(i).Set(reflect.Zero(valObj.Elem().Field(i).Type()))
				}
				hasContent = ret || hasContent

			}
		case reflect.Struct:
			ret := recursiveFlatten((valObj.Elem().Field(i).Addr().Interface()))
			hasContent = ret || hasContent
		case reflect.Slice:
			if typeObj.Elem().Field(i).Type == reflect.TypeOf([]string{}) {
				del := false
				for _, elem := range (valObj.Elem().Field(i).Interface()).([]string) {
					if strings.EqualFold(elem, undef) {
						del = true
					}
				}
				if del {
					valObj.Elem().Field(i).SetLen(0)
				}
			}
			if valObj.Elem().Field(i).Len() > 0 {
				hasContent = true
			}
		case reflect.String:
			if strings.EqualFold(valObj.Elem().Field(i).String(), undef) {
				valObj.Elem().Field(i).SetString("")
			}
			if valObj.Elem().Field(i).String() != "" {
				hasContent = true
			}
		case reflect.Bool:
			val := valObj.Elem().Field(i).Interface().(bool)
			hasContent = hasContent || val
		default:
			switch valObj.Elem().Field(i).Type() {
			case reflect.TypeOf(net.IP{}):
				val := valObj.Elem().Field(i).Interface().(net.IP)
				if len(val) != 0 && !val.IsUnspecified() {
					hasContent = true
				}
			default:
			}
		}
		if !hasContent {
			valObj.Elem().Field(i).Set(reflect.Zero(valObj.Elem().Field(i).Type()))
		}
	}
	return
}

/*
Create a string slice, where every element represents a yaml entry, used for node/profile edit
in order to get a summary of all available elements
*/
func UnmarshalConf(obj interface{}, excludeList []string) (lines []string) {
	objType := reflect.TypeOf(obj)
	// now iterate of every field
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		if field.Tag.Get("comment") != "" {
			if ymlStr, ok := getYamlString(field, excludeList); ok {
				lines = append(lines, ymlStr...)
			}
		}
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct && field.Tag.Get("yaml") != "" {
			typeLine := field.Tag.Get("yaml")
			if len(strings.Split(typeLine, ",")) > 1 {
				typeLine = strings.Split(typeLine, ",")[0] + ":"
			}
			lines = append(lines, typeLine)
			nestedLine := UnmarshalConf(reflect.New(field.Type.Elem()).Elem().Interface(), excludeList)
			for _, ln := range nestedLine {
				lines = append(lines, "  "+ln)
			}
		} else if field.Type.Kind() == reflect.Map && field.Type.Elem().Kind() == reflect.Ptr {
			typeLine := field.Tag.Get("yaml")
			if len(strings.Split(typeLine, ",")) > 1 {
				typeLine = strings.Split(typeLine, ",")[0] + ":"
			}
			lines = append(lines, typeLine, "  element:")
			nestedLine := UnmarshalConf(reflect.New(field.Type.Elem().Elem()).Elem().Interface(), excludeList)
			for _, ln := range nestedLine {
				lines = append(lines, "    "+ln)
			}
		} else if field.Type.Kind() == reflect.Struct && field.Anonymous {
			nestedLine := UnmarshalConf(reflect.New(field.Type).Elem().Interface(), excludeList)
			lines = append(lines, nestedLine...)
		}
	}
	return lines
}

/*
Get the string of the yaml tag
*/
func getYamlString(myType reflect.StructField, excludeList []string) ([]string, bool) {
	ymlStr := myType.Tag.Get("yaml")
	if len(strings.Split(ymlStr, ",")) > 1 {
		ymlStr = strings.Split(ymlStr, ",")[0]
	}
	if util.InSlice(excludeList, ymlStr) {
		return []string{""}, false
	} else if myType.Tag.Get("comment") == "" && myType.Type.Kind() == reflect.String {
		return []string{""}, false
	}
	if myType.Type.Kind() == reflect.String {
		fieldType := myType.Tag.Get("type")
		if fieldType == "" {
			fieldType = "string"
		}
		ymlStr += ": " + fieldType
		return []string{ymlStr}, true
	} else if myType.Type == reflect.TypeOf([]string{}) {
		return []string{ymlStr + ":", "  - string"}, true
	} else if myType.Type == reflect.TypeOf(map[string]string{}) {
		return []string{ymlStr + ":", "  key: value"}, true
	} else if myType.Type.Kind() == reflect.Ptr {
		return []string{ymlStr + ":"}, true
	}
	return []string{ymlStr}, true
}

/*
Getters for unexported fields
*/

/*
Returns the id of the node
*/
func (node *Node) Id() string {
	return node.id
}

/*
Returns the id of the profile
*/
func (node *Profile) Id() string {
	return node.id
}

/*
Returns if the node is a valid in the database
*/
func (node *Node) Valid() bool {
	return node.valid
}

func (node *Node) updatePrimaryNetDev() {
	if netdev, ok := node.NetDevs[node.PrimaryNetDev]; ok {
		netdev.primary = true
	} else {
		keys := make([]string, 0)
		for k := range node.NetDevs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if len(keys) > 0 {
			wwlog.Debug("%s: no primary defined, sanitizing to: %s", node.id, keys[0])
			node.NetDevs[keys[0]].primary = true
			node.PrimaryNetDev = keys[0]
		}
	}
}

/*
Check if the netdev is the primary one
*/
func (dev *NetDev) Primary() bool {
	return dev.primary
}

// returns all negated elements which are marked with ! as prefix
// from a list
func negList(list []string) (ret []string) {
	for _, tok := range list {
		if strings.HasPrefix(tok, "~") {
			ret = append(ret, tok[1:])
		}
	}
	return
}

// clean a list from negated tokens
func cleanList(list []string) (ret []string) {
	neg := negList(list)
	for _, listTok := range list {
		notNegate := true
		for _, negTok := range neg {
			if listTok == negTok || listTok == "~"+negTok {
				notNegate = false
			}
		}
		if notNegate {
			ret = append(ret, listTok)
		}
	}
	return ret
}

/*
Return the ipv4 address and mask in CIDR format. Aimed for the use in
templates.
*/
func (netdev *NetDev) IpCIDR() string {
	if netdev.Ipaddr.IsUnspecified() || netdev.Netmask.IsUnspecified() {
		return ""
	}
	ipCIDR := net.IPNet{
		IP:   netdev.Ipaddr,
		Mask: net.IPMask(netdev.Netmask),
	}
	return ipCIDR.String()
}
