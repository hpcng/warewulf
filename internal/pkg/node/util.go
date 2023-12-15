package node

import (
	"errors"
	"net"
	"reflect"
	"strings"
)

/*
Gets a node by its hardware(mac) address
*/
func (config *NodeYaml) FindByHwaddr(hwa string) (NodeConf, error) {
	if _, err := net.ParseMAC(hwa); err != nil {
		return NodeConf{}, errors.New("invalid hardware address: " + hwa)
	}
	nodeList, _ := config.FindAllNodes()
	for _, node := range nodeList {
		for _, dev := range node.NetDevs {
			if strings.EqualFold(dev.Hwaddr, hwa) {
				return node, nil
			}
		}
	}

	return NodeConf{}, ErrNotFound
}

/*
Find a node by its ip address
*/
func (config *NodeYaml) FindByIpaddr(ipaddr string) (NodeConf, error) {
	addr := net.ParseIP(ipaddr)
	if addr == nil {
		return NodeConf{}, errors.New("invalid IP:" + ipaddr)
	}
	nodeList, err := config.FindAllNodes()
	if err != nil {
		return NodeConf{}, err
	}
	for _, node := range nodeList {
		for _, dev := range node.NetDevs {
			if dev.Ipaddr.Equal(addr) {
				return node, nil
			}
		}
	}

	return NodeConf{}, ErrNotFound
}

/*
Check if the Object is empty, has no valid values
*/
func ObjectIsEmpty(obj interface{}) bool {
	if obj == nil {
		return true
	}
	varType := reflect.TypeOf(obj)
	varVal := reflect.ValueOf(obj)
	if varType.Kind() == reflect.Ptr && !varVal.IsNil() {
		return ObjectIsEmpty(varVal.Elem().Interface())
	}
	if varVal.IsZero() {
		return true
	}
	for i := 0; i < varType.NumField(); i++ {
		if varType.Field(i).Type.Kind() == reflect.String && !varVal.Field(i).IsZero() {
			val := varVal.Field(i).Interface().(string)
			if val != "" {
				return false
			}
		} else if varType.Field(i).Type == reflect.TypeOf(map[string]string{}) {
			if len(varVal.Field(i).Interface().(map[string]string)) != 0 {
				return false
			}
		} else if varType.Field(i).Type.Kind() == reflect.Ptr {
			if !ObjectIsEmpty(varVal.Field(i).Interface()) {
				return false
			}
		} else if varType.Field(i).Type == reflect.TypeOf(net.IP{}) {
			val := varVal.Field(i).Interface().(net.IP)
			if len(val) != 0 && !val.IsUnspecified() {
				return false
			}
		} else if varType.Field(i).Type == reflect.TypeOf(net.IPMask{}) {
			o, b := varVal.Field(i).Interface().(net.IPMask).Size()
			if o == 0 && b == 0 {
				return false
			}
		}

	}
	return true
}
