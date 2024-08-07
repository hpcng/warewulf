package node

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

type sortByName []NodeInfo

func (a sortByName) Len() int           { return len(a) }
func (a sortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByName) Less(i, j int) bool { return a[i].Id.Print() < a[j].Id.Print() }

func GetUnsetVerbs() []string {
	return []string{"UNSET", "DELETE", "UNDEF", "undef", "--", "nil", "0.0.0.0"}
}

/**********
 *
 * Filters
 *
 *********/

/*
Filter a given slice of NodeInfo against a given
regular expression
*/
func FilterByName(set []NodeInfo, searchList []string) []NodeInfo {
	var ret []NodeInfo
	unique := make(map[string]NodeInfo)

	if len(searchList) > 0 {
		for _, search := range searchList {
			for _, entry := range set {
				if match, _ := regexp.MatchString("^"+search+"$", entry.Id.Get()); match {
					unique[entry.Id.Get()] = entry
				}
			}
		}
		for _, n := range unique {
			ret = append(ret, n)
		}
	} else {
		ret = set
	}

	sort.Sort(sortByName(ret))
	return ret
}

/*
Filter a given map of NodeConf against given regular expression.
*/
func FilterMapByName(inputMap map[string]*NodeConf, searchList []string) (retMap map[string]*NodeConf) {
	retMap = map[string]*NodeConf{}
	if len(searchList) > 0 {
		for _, search := range searchList {
			for name, nConf := range inputMap {
				if match, _ := regexp.MatchString("^"+search+"$", name); match {
					retMap[name] = nConf
				}
			}
		}
	}
	return retMap
}

/**********
 *
 * Sets
 *
 *********/

/*
Set value. If argument is 'UNDEF', 'DELETE',
'UNSET" or '--' the value is removed.
N.B. the '--' might never ever happen as '--'
is parsed out by cobra
*/
func (ent *Entry) Set(val string) {
	if val == "" {
		return
	}

	if util.InSlice(GetUnsetVerbs(), val) {
		wwlog.Debug("Removing value for %v", *ent)
		ent.value = []string{""}
	} else {
		ent.value = []string{val}
	}
}

/*
Set bool
*/
func (ent *Entry) SetB(val bool) {
	if val {
		ent.value = []string{"true"}
	} else {
		ent.value = []string{"false"}
	}
}

func (ent *Entry) SetSlice(val []string) {
	if len(val) == 0 {
		return
	} else if len(val) == 1 && val[0] == "" { // check also for an "empty" slice
		return
	}
	ent.isSlice = true
	if util.InSlice(GetUnsetVerbs(), val[0]) {
		ent.value = []string{}
	} else {
		ent.value = val
	}
}

/*
Set alternative value
*/
func (ent *Entry) SetAlt(val string, from string) {
	if val == "" {
		return
	}
	ent.altvalue = []string{val}
	ent.from = from
}

/*
Sets alternative bool
*/
func (ent *Entry) SetAltB(val bool, from string) {
	if val {
		ent.altvalue = []string{"true"}
		ent.from = from
	} else {
		ent.altvalue = []string{"false"}
		ent.from = from
	}
}

/*
Sets alternative slice
*/
func (ent *Entry) SetAltSlice(val []string, from string) {
	if len(val) == 0 {
		return
	}
	ent.isSlice = true
	ent.altvalue = append(ent.altvalue, val...)
	if ent.from == "" {
		ent.from = from
	} else {
		ent.from = ent.from + "," + from
	}
}

/*
Sets the default value of an entry.
*/
func (ent *Entry) SetDefault(val string) {
	if val == "" {
		return
	}
	ent.def = []string{val}
}

/*
Set the default entry as slice
*/
func (ent *Entry) SetDefaultSlice(val []string) {
	if len(val) == 0 {
		return
	}
	ent.isSlice = true
	ent.def = val
}

/*
Set default entry as bool
*/
func (ent *Entry) SetDefaultB(val bool) {
	if val {
		ent.def = []string{"true"}
	} else {
		ent.def = []string{"false"}
	}
}

/*
Remove a element from a slice
*/
func (ent *Entry) SliceRemoveElement(val string) {
	util.SliceRemoveElement(ent.value, val)
}

/**********
*
* Gets
*
*********/
/*
Gets the the entry of the value in following order
* node value if set
* profile value if set
* default value if set
*/
func (ent *Entry) Get() string {
	if len(ent.value) != 0 {
		return ent.value[0]
	}
	if len(ent.altvalue) != 0 {
		return ent.altvalue[0]
	}
	if len(ent.def) != 0 {
		return ent.def[0]
	}
	return ""
}

/*
Get the bool value of an entry.
*/
func (ent *Entry) GetB() bool {
	if len(ent.value) > 0 {
		return !(strings.ToLower(ent.value[0]) == "false" ||
			strings.ToLower(ent.value[0]) == "no" ||
			ent.value[0] == "0")
	} else if len(ent.altvalue) > 0 {
		return !(strings.ToLower(ent.altvalue[0]) == "false" ||
			strings.ToLower(ent.altvalue[0]) == "no" ||
			ent.altvalue[0] == "0")
	} else {
		return !(len(ent.def) == 0 ||
			strings.ToLower(ent.def[0]) == "false" ||
			strings.ToLower(ent.def[0]) == "no" ||
			ent.def[0] == "0")
	}
}

// returns all negated elemets which are marked with ! as prefix
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
Returns a string slice created from a comma separated list of the value.
*/
func (ent *Entry) GetSlice() []string {
	retval := append(ent.value, ent.altvalue...)
	if len(retval) != 0 {
		return cleanList(retval)
	}
	if len(ent.def) != 0 {
		return ent.def
	}
	return []string{}
}

func (ent *Entry) GetSlicePreferAlt() []string {
	retval := append(ent.altvalue, ent.value...)
	if len(retval) != 0 {
		return cleanList(retval)
	}
	if len(ent.def) != 0 {
		return ent.def
	}
	return []string{}
}

/*
Get the real value, not the alternative of default one.
*/
func (ent *Entry) GetReal() string {
	if len(ent.value) == 0 {
		return ""
	}
	return ent.value[0]
}

/*
Get the real value, not the alternative of default one.
*/
func (ent *Entry) GetRealSlice() []string {
	if len(ent.value) == 0 {
		return []string{}
	}
	return ent.value
}

/*
true if the entry has set a real value, else false.
*/
func (ent *Entry) GotReal() bool {
	return len(ent.value) != 0
}

/*
Get a pointer to the value
*/
func (ent *Entry) GetPointer() *string {
	ret := ent.Get()
	return &ret
}

/*
Try to get a int of a value, 0 if value can't be parsed!
*/
func (ent *Entry) GetInt() int {
	var ret int
	if len(ent.value) != 0 {
		ret, _ = strconv.Atoi(ent.value[0])
	} else if len(ent.altvalue) != 0 {
		ret, _ = strconv.Atoi(ent.altvalue[0])
	} else if len(ent.def) != 0 {
		ret, _ = strconv.Atoi(ent.def[0])
	}
	return ret
}

/*
Ptr to int
*/
func (ent *Entry) GetIntPtr() *int {
	ret := ent.GetInt()
	return &ret
}

/**********
 *
 * Misc
 *
 *********/
/*
Gets the the entry of the value in following order
* node value if set
* profile value if set
* default value if set
*/
func (ent *Entry) Print() string {
	if !ent.isSlice {
		if len(ent.value) != 0 {
			return ent.value[0]
		}
		if len(ent.altvalue) != 0 {
			return ent.altvalue[0]
		}
		if len(ent.def) != 0 {
			return "(" + ent.def[0] + ")"
		}
	} else {
		var ret string
		if len(ent.value) != 0 || len(ent.altvalue) != 0 {
			combList := append(ent.value, ent.altvalue...)
			ret = strings.Join(cleanList(combList), ",")
			if len(negList(combList)) > 0 {
				ret += " ~{" + strings.Join(negList(combList), ",") + "}"
			}

		}
		if ret != "" {
			return ret
		}
		if len(ent.def) != 0 {
			return "(" + strings.Join(ent.def, ",") + ")"
		}

	}
	return "--"
}

/* Gets the the entry of the value in following order
* profile value if set
* node value if set
* default value if set
 */
func (ent *Entry) PrintPreferAlt() string {
	if !ent.isSlice {
		if len(ent.altvalue) != 0 {
			return ent.altvalue[0]
		}
		if len(ent.value) != 0 {
			return ent.value[0]
		}
		if len(ent.def) != 0 {
			return "(" + ent.def[0] + ")"
		}
	} else {
		var ret string
		if len(ent.value) != 0 || len(ent.altvalue) != 0 {
			combList := append(ent.altvalue, ent.value...)
			ret = strings.Join(cleanList(combList), ",")
			if len(negList(combList)) > 0 {
				ret += " ~{" + strings.Join(negList(combList), ",") + "}"
			}

		}
		if ret != "" {
			return ret
		}
		if len(ent.def) != 0 {
			return "(" + strings.Join(ent.def, ",") + ")"
		}

	}
	return "--"
}

/*
Was used for combined stringSlice

func (ent *Entry) PrintComb() string {
	if ent.value != "" && ent.altvalue != "" {
		return "[" + ent.value + "," + ent.altvalue + "]"
	}
	return ent.Print()
}
*/

/*
same as GetB()
*/
func (ent *Entry) PrintB() string {
	if len(ent.value) != 0 || len(ent.altvalue) != 0 {
		return fmt.Sprintf("%t", ent.GetB())
	}
	return fmt.Sprintf("(%t)", ent.GetB())
}

/*
Returns SUPERSEDED if value was set per node or
per profile. Else -- is returned.
*/
func (ent *Entry) Source() string {
	if len(ent.value) != 0 && len(ent.altvalue) != 0 {
		return "SUPERSEDED"
		// return fmt.Sprintf("[%s]", ent.from)
	} else if ent.from == "" {
		return "--"
	}
	return ent.from
}

/*
Check if value was defined.
*/
func (ent *Entry) Defined() bool {
	if len(ent.value) != 0 {
		return true
	}
	if len(ent.altvalue) != 0 {
		return true
	}
	if len(ent.def) != 0 {
		return true
	}
	return false
}

/*
Create an empty node NodeConf
*/
func NewConf() (nodeconf NodeConf) {
	nodeconf.Ipmi = new(IpmiConf)
	nodeconf.Ipmi.Tags = map[string]string{}
	nodeconf.Kernel = new(KernelConf)
	nodeconf.NetDevs = make(map[string]*NetDevs)
	nodeconf.Tags = map[string]string{}
	return nodeconf
}

/*
Create an empty node NodeInfo
*/
func NewInfo() (nodeInfo NodeInfo) {
	nodeInfo.Ipmi = new(IpmiEntry)
	nodeInfo.Ipmi.Tags = map[string]*Entry{}
	nodeInfo.Kernel = new(KernelEntry)
	nodeInfo.NetDevs = make(map[string]*NetDevEntry)
	nodeInfo.Tags = make(map[string]*Entry)
	return nodeInfo
}

/*
Get a entry by its name
*/
func GetByName(node interface{}, name string) (string, error) {
	valEntry := reflect.ValueOf(node)
	entryField := valEntry.Elem().FieldByName(name)
	if entryField == (reflect.Value{}) {
		return "", fmt.Errorf("couldn't find field with name: %s", name)
	}
	if entryField.Type() != reflect.TypeOf(Entry{}) {
		return "", fmt.Errorf("field %s is not of type node.Entry", name)
	}
	myEntry := entryField.Interface().(Entry)
	return myEntry.Get(), nil
}

/*
Check if the Netdev is empty, so has no values set
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
		}
	}
	return true
}
