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
func (config *NodesYaml) FindByHwaddr(hwa string) (Node, error) {
	if _, err := net.ParseMAC(hwa); err != nil {
		return Node{}, errors.New("invalid hardware address: " + hwa)
	}
	nodeList, _ := config.FindAllNodes()
	for _, node := range nodeList {
		for _, dev := range node.NetDevs {
			if strings.EqualFold(dev.Hwaddr, hwa) {
				return node, nil
			}
		}
	}

	return Node{}, ErrNotFound
}

/*
Find a node by its ip address
*/
func (config *NodesYaml) FindByIpaddr(ipaddr string) (Node, error) {
	addr := net.ParseIP(ipaddr)
	if addr == nil {
		return Node{}, errors.New("invalid IP:" + ipaddr)
	}
	nodeList, err := config.FindAllNodes()
	if err != nil {
		return Node{}, err
	}
	for _, node := range nodeList {
		for _, dev := range node.NetDevs {
			if dev.Ipaddr.Equal(addr) {
				return node, nil
			}
		}
	}

	return Node{}, ErrNotFound
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
			val := varVal.Field(i).String()
			if val != "" {
				return false
			}
		} else if varType.Field(i).Type == reflect.TypeOf(map[string]string{}) {
			if varVal.Field(i).Len() != 0 {
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
		}
	}
	return true
}
