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
	Net         string            `lopt:"netname" comment:"network which is modified" default:"default"`
	DiskName    string            `lopt:"diskname" comment:"set diskdevice name"`
	PartName    string            `lopt:"partname" comment:"set the partition name so it can be used by a file system"`
	FsName      string            `lopt:"fsname" comment:"set the file system name which must match a partition name"`
}

/*
Create cmd line flags from the NodeConf fields. Returns a []func() where every function must be called, as the command line parser returns e.g. netip.IP objects which must be parsed
back to strings.
*/
func (nodeConf *Node) CreateFlags(baseCmd *cobra.Command) {
	recursiveCreateFlags(nodeConf, baseCmd)
}

func (profileConf *Profile) CreateFlags(baseCmd *cobra.Command) {
	recursiveCreateFlags(profileConf, baseCmd)
}

func (del *NodeConfDel) CreateDelFlags(baseCmd *cobra.Command) {
	recursiveCreateFlags(del, baseCmd)

}
func (add *NodeConfAdd) CreateAddFlags(baseCmd *cobra.Command) {
	recursiveCreateFlags(add, baseCmd)

}

func recursiveCreateFlags(obj interface{}, baseCmd *cobra.Command) {
	elemType := reflect.TypeOf(obj).Elem()
	elemVal := reflect.ValueOf(obj).Elem()

	for i := 0; i < elemVal.NumField(); i++ {
		field := elemType.Field(i)
		fieldVal := elemVal.Field(i)

		if !field.IsExported() {
			continue
		}

		if field.Tag.Get("comment") != "" {
			createFlags(baseCmd, field, &fieldVal)

		} else if field.Anonymous {
			recursiveCreateFlags(fieldVal.Addr().Interface(), baseCmd)

		} else if field.Type.Kind() == reflect.Ptr {
			recursiveCreateFlags(fieldVal.Interface(), baseCmd)

		} else if field.Type.Kind() == reflect.Struct {
			recursiveCreateFlags(fieldVal.Addr().Interface(), baseCmd)

		} else if field.Type.Kind() == reflect.Map {
			switch field.Type.Elem().Kind() {
			case reflect.String, reflect.Interface:
				continue
			case reflect.Pointer, reflect.Slice, reflect.Map:
				// add a map with key UNDEF so that it can hold values N.B. UNDEF can never be added through command line
				key := reflect.ValueOf("UNDEF")
				if fieldVal.Len() == 0 {
					if fieldVal.IsNil() {
						fieldVal.Set(reflect.MakeMap(field.Type))
					}
					newPtr := reflect.New(field.Type.Elem().Elem())
					fieldVal.SetMapIndex(key, newPtr)
				} else {
					key = fieldVal.MapKeys()[0]
				}
				recursiveCreateFlags(fieldVal.MapIndex(key).Interface(), baseCmd)
			}
		}
	}
}

/*
Helper function to create the different PersistentFlags() for different types.
*/
func createFlags(baseCmd *cobra.Command,
	myType reflect.StructField, myVal *reflect.Value) {
	var wwbool wwtype.WWbool
	if myType.Tag.Get("lopt") != "" {
		if myType.Type == reflect.TypeOf("") {
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
		} else if myType.Type == reflect.TypeOf(wwbool) {
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
