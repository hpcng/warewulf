package node

import (
	"regexp"
	"strings"
)

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
				b, _ := regexp.MatchString("^"+search+"$", entry.Id.Get())
				if b {
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

	return ret
}

/**********
 *
 * Sets
 *
 *********/

/*
 Set value. If argument is 'UNDEF', 'DELETE',
 'UNSET" or '--'. The value is removed.
 N.B. the '--' might never ever happen as '--'
 is parsed out by cobra
*/
func (ent *Entry) Set(val string) {
	if val == "" {
		return
	}

	if val == "UNDEF" || val == "DELETE" || val == "UNSET" || val == "--" {
		ent.value = []string{}
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
	}
}

func (ent *Entry) SetSlice(val []string) {
	if len(val) == 0 {
		return
	}
	if val[0] == "UNDEF" || val[0] == "DELETE" || val[0] == "UNSET" || val[0] == "--" {
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
	}
}

/*
Sets alternative slice
*/
func (ent *Entry) SetAltSlice(val []string, from string) {
	if len(val) == 0 {
		return
	}
	ent.altvalue = val
	ent.from = from
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
	ent.def = val

}

/**********
*
* Gets
*
*********/
/*
Gets the the entry of the value in folowing order
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
	if len(ent.value) == 0 || ent.value[0] == "false" || ent.value[0] == "no" {
		if len(ent.altvalue) == 0 || ent.altvalue[0] == "false" || ent.altvalue[0] == "no" {
			return false
		}
		return false
	}
	return true
}

/*
Returns a string slice created from a comma seperated list of the value.
*/
func (ent *Entry) GetSlice() []string {
	var retval []string
	if len(ent.value) != 0 {
		return ent.value
	}
	if len(ent.altvalue) != 0 {
		return ent.altvalue
	}
	if len(ent.def) != 0 {
		return ent.def
	}
	return retval
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

/**********
 *
 * Misc
 *
 *********/

/*
Returns the value of Entry if it was defined set or
alternative is presend. Default value is in '()'. If
nothing is defined '--' is returned.
*/
func (ent *Entry) Print() string {
	if len(ent.value) != 0 {
		return strings.Join(ent.value, ",")
	}
	if len(ent.altvalue) != 0 {
		return strings.Join(ent.altvalue, ",")
	}
	if len(ent.def) != 0 {
		return "(" + strings.Join(ent.def, ",") + ")"
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
func (ent *Entry) PrintB() bool {
	return ent.GetB()
}

/*
Returns SUPERSEDED if value was set per node or
per profile. Else -- is returned.
*/
func (ent *Entry) Source() string {
	if len(ent.value) != 0 && len(ent.altvalue) != 0 {
		return "SUPERSEDED"
		//return fmt.Sprintf("[%s]", ent.from)
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
