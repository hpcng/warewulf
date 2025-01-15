package node

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwtype"
)

/*
Checks if for NodeConf all values can be parsed according to their type.
*/
func (nodeConf *Node) Check() (err error) {
	nodeInfoType := reflect.TypeOf(nodeConf)
	nodeInfoVal := reflect.ValueOf(nodeConf)
	return check(nodeInfoType, nodeInfoVal)
}

func (profileConf *Profile) Check() (err error) {
	profileInfoType := reflect.TypeOf(profileConf)
	profileInfoVal := reflect.ValueOf(profileConf)
	return check(profileInfoType, profileInfoVal)
}

func check(infoType reflect.Type, infoVal reflect.Value) (err error) {
	if !infoVal.Elem().IsValid() {
		return nil
	}
	// now iterate of every field
	for i := 0; i < infoVal.Elem().NumField(); i++ {
		if !infoType.Elem().Field(i).IsExported() {
			continue
		} else if infoType.Elem().Field(i).Type.Kind() == reflect.String {
			newFmt, err := checker(infoVal.Elem().Field(i).String(), infoType.Elem().Field(i).Tag.Get("type"))
			if err != nil {
				return fmt.Errorf("field: %s value:%s err: %s", infoType.Elem().Field(i).Name, infoVal.Elem().Field(i).String(), err)
			} else if newFmt != "" {
				infoVal.Elem().Field(i).SetString(newFmt)
			}
		} else if infoType.Elem().Field(i).Type.Kind() == reflect.Ptr && !infoVal.Elem().Field(i).IsNil() {
			nestType := reflect.TypeOf(infoVal.Elem().Field(i).Interface())
			nestVal := reflect.ValueOf(infoVal.Elem().Field(i).Interface())
			for j := 0; j < nestType.Elem().NumField(); j++ {
				if nestType.Elem().Field(j).Type.Kind() == reflect.String {
					newFmt, err := checker(nestVal.Elem().Field(j).Interface().(string), nestType.Elem().Field(j).Tag.Get("type"))
					if err != nil {
						return fmt.Errorf("field: %s value:%s err: %s", nestType.Elem().Field(j).Name, nestVal.Elem().Field(j).String(), err)
					} else if newFmt != "" {
						nestVal.Elem().Field(j).SetString(newFmt)
					}
				}
			}
		} else if infoType.Elem().Field(i).Type == reflect.TypeOf(map[string]*NetDev(nil)) {
			netMap := infoVal.Elem().Field(i).Interface().(map[string]*NetDev)
			for _, val := range netMap {
				netType := reflect.TypeOf(val)
				netVal := reflect.ValueOf(val)
				for j := 0; j < netType.Elem().NumField(); j++ {
					newFmt, err := checker(netVal.Elem().Field(j).String(), netType.Elem().Field(j).Tag.Get("type"))
					if err != nil {
						return fmt.Errorf("field: %s value:%s err: %s", netType.Elem().Field(j).Name, netVal.Elem().Field(j).String(), err)
					} else if newFmt != "" {
						netVal.Elem().Field(j).SetString(newFmt)
					}
				}
			}
		}
	}
	return nil
}

func checker(value string, valType string) (niceValue string, err error) {
	if valType == "" || value == "" || util.InSlice(wwtype.GetUnsetVerbs(), value) {
		return "", nil
	}
	switch valType {
	case "":
		return "", nil
	case "bool":
		if strings.ToLower(value) == "yes" {
			return "true", nil
		}
		if strings.ToLower(value) == "no" {
			return "false", nil
		}
		myBool, err := strconv.ParseBool(value)
		return strconv.FormatBool(myBool), err
	case "IP":
		if addr := net.ParseIP(value); addr == nil {
			return "", fmt.Errorf("%s can't be parsed to ip address", value)
		} else {
			return addr.String(), nil
		}
	case "MAC":
		if mac, err := net.ParseMAC(value); err != nil {
			return "", fmt.Errorf("%s can't be parsed to MAC address: %s", value, err)
		} else {
			return mac.String(), nil
		}
	case "uint":
		if _, err := strconv.ParseUint(value, 10, 64); err != nil {
			return "", fmt.Errorf("%s is not a uint: %s", value, err)
		}
	}
	return "", nil
}
