package node

import (
	"net"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/pkg/wwtype"
)

type NodeConfDel struct {
	TagsDel     []string `lopt:"tagdel" comment:"add tags"`
	IpmiTagsDel []string `lopt:"ipmitagdel" comment:"delete ipmi tags"`
	NetTagsDel  []string `lopt:"nettagdel" comment:"delete network tags"`
	NetDel      string   `lopt:"netdel" comment:"network to delete"`
	DiskDel     string   `lopt:"diskdel" comment:"delete the disk from the configuration"`
	PartDel     string   `lopt:"partdel" comment:"delete the partition from the configuration"`
	FsDel       string   `lopt:"fsdel" comment:"delete the fs from the configuration"`
}
type NodeConfAdd struct {
	TagsAdd     map[string]string `lopt:"tagadd" comment:"add tags"`
	IpmiTagsAdd map[string]string `lopt:"ipmitagadd" comment:"add ipmi tags"`
	NetTagsAdd  map[string]string `lopt:"nettagadd" comment:"add network tags"`
	Net         string            `lopt:"netname" comment:"network which is modified"`
	DiskName    string            `lopt:"diskname" comment:"set diskdevice name"`
	PartName    string            `lopt:"partname" comment:"set the partition name so it can be used by a file system"`
	FsName      string            `lopt:"fsname" comment:"set the file system name which must match a partition name"`
}

/*
Create cmd line flags from the NodeConf fields. Returns a []func() where every function must be called, as the command line parser returns e.g. netip.IP objects which must be parsed
back to strings.
*/
func (nodeConf *NodeConf) CreateFlags(baseCmd *cobra.Command) {
	recursiveCreateFlags(nodeConf, baseCmd)
}

func (del *NodeConfDel) CreateDelFlags(baseCmd *cobra.Command) {
	recursiveCreateFlags(del, baseCmd)

}
func (add *NodeConfAdd) CreateAddFlags(baseCmd *cobra.Command) {
	recursiveCreateFlags(add, baseCmd)

}

func recursiveCreateFlags(obj interface{}, baseCmd *cobra.Command) {
	// now iterate of every field
	nodeInfoType := reflect.TypeOf(obj)
	nodeInfoVal := reflect.ValueOf(obj)
	for i := 0; i < nodeInfoVal.Elem().NumField(); i++ {
		if nodeInfoType.Elem().Field(i).Tag.Get("comment") != "" {
			field := nodeInfoVal.Elem().Field(i)
			createFlags(baseCmd, nodeInfoType.Elem().Field(i), &field)

		} else if nodeInfoType.Elem().Field(i).Type.Kind() == reflect.Ptr {
			recursiveCreateFlags(nodeInfoVal.Elem().Field(i).Interface(), baseCmd)

		} else if nodeInfoType.Elem().Field(i).Type.Kind() == reflect.Map &&
			nodeInfoType.Elem().Field(i).Type != reflect.TypeOf(map[string]string{}) {
			// add a map with key UNDEF so that it can hold values N.B. UNDEF can never be added through command line
			key := reflect.ValueOf("UNDEF")
			if nodeInfoVal.Elem().Field(i).Len() == 0 {
				if nodeInfoVal.Elem().Field(i).IsNil() {
					nodeInfoVal.Elem().Field(i).Set(reflect.MakeMap(nodeInfoType.Elem().Field(i).Type))
				}
				newPtr := reflect.New(nodeInfoType.Elem().Field(i).Type.Elem().Elem())
				nodeInfoVal.Elem().Field(i).SetMapIndex(key, newPtr)
			} else {
				key = nodeInfoVal.Elem().Field(i).MapKeys()[0]
			}
			recursiveCreateFlags(nodeInfoVal.Elem().Field(i).MapIndex(key).Interface(), baseCmd)
		} else if nodeInfoType.Elem().Field(i).Anonymous {
			recursiveCreateFlags(nodeInfoVal.Elem().Field(i).Addr().Interface(), baseCmd)
		} else if nodeInfoType.Elem().Field(i).Type.Kind() == reflect.Struct {
			recursiveCreateFlags(nodeInfoVal.Elem().Field(i).Addr().Interface(), baseCmd)
		}
	}
}

/*
Helper function to create the different PersistentFlags() for different types.
*/
func createFlags(baseCmd *cobra.Command,
	myType reflect.StructField, myVal *reflect.Value) {
	if myType.Tag.Get("lopt") != "" {
		if myType.Type.Kind() == reflect.String {
			ptr := myVal.Addr().Interface().(*string)
			if myType.Tag.Get("sopt") != "" {
				baseCmd.PersistentFlags().StringVarP(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("sopt"),
					myType.Tag.Get("default"),
					myType.Tag.Get("comment"))
			} else {
				baseCmd.PersistentFlags().StringVar(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("default"),
					myType.Tag.Get("comment"))
			}

		} else if myType.Type == reflect.TypeOf([]string{}) {
			ptr := myVal.Addr().Interface().(*[]string)
			if myType.Tag.Get("sopt") != "" {
				baseCmd.PersistentFlags().StringSliceVarP(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("sopt"),
					[]string{},
					myType.Tag.Get("comment"))
			} else {
				baseCmd.PersistentFlags().StringSliceVar(ptr,
					myType.Tag.Get("lopt"),
					[]string{},
					myType.Tag.Get("comment"))

			}
		} else if myType.Type == reflect.TypeOf(map[string]string{}) {
			ptr := myVal.Addr().Interface().(*map[string]string)
			if myType.Tag.Get("sopt") != "" {
				baseCmd.PersistentFlags().StringToStringVarP(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("sopt"),
					map[string]string{}, // empty default!
					myType.Tag.Get("comment"))
			} else {
				baseCmd.PersistentFlags().StringToStringVar(ptr,
					myType.Tag.Get("lopt"),
					map[string]string{}, // empty default!
					myType.Tag.Get("comment"))
			}
		} else if myType.Type == reflect.TypeOf(true) {
			ptr := myVal.Addr().Interface().(*bool)
			if myType.Tag.Get("sopt") != "" {
				baseCmd.PersistentFlags().BoolVarP(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("sopt"),
					false, // empty default!
					myType.Tag.Get("comment"))
			} else {
				baseCmd.PersistentFlags().BoolVar(ptr,
					myType.Tag.Get("lopt"),
					false, // empty default!
					myType.Tag.Get("comment"))
			}

		} else if myType.Type == reflect.TypeOf(true) {
			ptr := myVal.Addr().Interface().(*bool)
			if myType.Tag.Get("sopt") != "" {
				baseCmd.PersistentFlags().BoolVarP(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("sopt"),
					false, // empty default!
					myType.Tag.Get("comment"))
			} else {
				baseCmd.PersistentFlags().BoolVar(ptr,
					myType.Tag.Get("lopt"),
					false, // empty default!
					myType.Tag.Get("comment"))
			}
		} else if myType.Type == reflect.TypeOf(net.IP{}) {
			ptr := myVal.Addr().Interface().(*net.IP)
			if myType.Tag.Get("sopt") != "" {
				baseCmd.PersistentFlags().IPVarP(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("sopt"),
					net.IP{}, // empty default!
					myType.Tag.Get("comment"))
			} else {
				baseCmd.PersistentFlags().IPVar(ptr,
					myType.Tag.Get("lopt"),
					net.IP{}, // empty default!
					myType.Tag.Get("comment"))
			}
		} else if myType.Type == reflect.TypeOf(net.IPMask{}) {
			ptr := myVal.Addr().Interface().(*net.IPMask)
			if myType.Tag.Get("sopt") != "" {
				baseCmd.PersistentFlags().IPMaskVarP(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("sopt"),
					net.IPMask{}, // empty default!
					myType.Tag.Get("comment"))
			} else {
				baseCmd.PersistentFlags().IPMaskVar(ptr,
					myType.Tag.Get("lopt"),
					net.IPMask{}, // empty default!
					myType.Tag.Get("comment"))
			}
		} else if myType.Type == reflect.TypeOf(wwtype.WWbool{}) {
			ptr := myVal.Addr().Interface().(*wwtype.WWbool)
			if myType.Tag.Get("sopt") != "" {
				baseCmd.PersistentFlags().VarP(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("sopt"),
					myType.Tag.Get("comment"))
				baseCmd.Flag(myType.Tag.Get("lopt")).NoOptDefVal = "true"
			} else {
				baseCmd.PersistentFlags().Var(ptr,
					myType.Tag.Get("lopt"),
					myType.Tag.Get("comment"))
				baseCmd.Flag(myType.Tag.Get("lopt")).NoOptDefVal = "true"
			}
		}
	}
}
