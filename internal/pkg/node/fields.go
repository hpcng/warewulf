package node

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// Field represents a single attribute field (typically to be used with a Node or a Profile),
// including the field's name, the source of the value (such as a node or profile ID), and
// the actual string value.
//
// It is primarily used to provide desired output for `wwctl <node|profile> list -a`.
type Field struct {
	Field  string
	Source string
	Value  string
}

// Set updates the field with the given source and value. If the value is empty, the operation
// is skipped. If the field already has a source, and the new source is empty, the previous source
// is marked as "SUPERSEDED" to indicate it was overridden (typically by a node's local
// configuration).
func (f *Field) Set(src, val string) {
	if val == "" {
		return
	}
	f.Value = val

	if f.Source != "" && src == "" {
		f.Source = "SUPERSEDED"
	} else {
		f.Source = src
	}
}

// fieldMap maps field names to Field objects. This structure is used to track and manage
// multiple fields, along with their sources and values, particularly by MergeNode.
type fieldMap map[string]*Field

// Set updates the correct field in the fieldMap with the given source and value.
// If the field does not already exist in the fieldMap, it is created.
func (fields fieldMap) Set(name, source, value string) {
	if fields[name] == nil {
		fields[name] = &Field{Field: name}
	}
	fields[name].Set(source, value)
}

// Source returns the source of the given field name if it exists in the map. If the field does
// not exist, an empty string is returned.
func (fields fieldMap) Source(name string) string {
	if field, ok := fields[name]; ok {
		return field.Source
	}
	return ""
}

// Value returns the value of the given field name if it exists in the map. If the field does
// not exist, an empty string is returned.
func (fields fieldMap) Value(name string) string {
	if field, ok := fields[name]; ok {
		return field.Value
	}
	return ""
}

// List returns a slice of Field structs for all fields that exist in the fieldMap, in the
// order they are defined on the provided object. This ensures a consistent ordering of fields
// for display purposes.
func (fields fieldMap) List(obj interface{}) (output []Field) {
	for _, name := range listFields(obj) {
		if field, ok := fields[name]; ok {
			output = append(output, *field)
		}
	}
	return output
}

// GetFieldList extracts all fields from the provided object and returns them as a slice of Fields.
// Each Field includes the field name and its string value. Fields that cannot be retrieved
// or converted are skipped.
func GetFieldList(obj interface{}) (fields []Field) {
	for _, name := range listFields(obj) {
		if value, err := getNestedFieldString(obj, name); err == nil {
			fields = append(fields, Field{Field: name, Value: value})
		}
	}
	return fields
}

var mapFieldElement *regexp.Regexp

func init() {
	// mapFieldElement matches map-indexed fields like "FieldName[Key]" to split into (FieldName, Key).
	mapFieldElement = regexp.MustCompile(`^([^[]+)\[([^\]]+)\]$`)
}

// getNestedFieldValue retrieves the reflect.Value of a nested field.
//
// Supported syntax:
// - Struct fields identified by a dotted path name (Struct.field)
// - Map keys identified by square brackets (Map[key])
//
// Pointers are automatically dereferenced.
//
// If any element of the path does not exist, an error is returned.
func getNestedFieldValue(obj interface{}, name string) (value reflect.Value, err error) {
	value = reflect.ValueOf(obj)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	fieldNames := strings.Split(name, ".")
	for _, fieldName := range fieldNames {
		var key string
		fieldName, key = parseMapField(fieldName)
		if value.Kind() == reflect.Pointer {
			if value.IsNil() {
				err = fmt.Errorf("no value: %v", name)
				return
			}
			value = value.Elem()
		}
		if !value.IsValid() {
			err = fmt.Errorf("no value: %v", name)
			return
		}
		value = value.FieldByName(fieldName)
		if key != "" {
			value = value.MapIndex(reflect.ValueOf(key))
			if !value.IsValid() {
				err = fmt.Errorf("no value: %v", name)
				return
			}
		}
	}
	return
}

// getNestedFieldString retrieves the string representation
// of a nested field as returned by getNestedFieldValue.
//
// Returns an error if the field does not exist or cannot be retrieved.
func getNestedFieldString(obj interface{}, name string) (string, error) {
	if value, err := getNestedFieldValue(obj, name); err != nil {
		return "", err
	} else {
		return valueStr(value), nil
	}
}

// parseMapField extracts the map key if the field name represents a map access (e.g. "Fields[key]" returns "Fields", "key").
// If there is no key specified, it simply returns the field name as is.
func parseMapField(name string) (field, key string) {
	if matches := mapFieldElement.FindStringSubmatch(name); matches != nil {
		return matches[1], matches[2]
	}
	return name, ""
}

// listFields returns a slice of strings representing all exported, visible fields of the given
// object's type, including nested fields in structs and keys in maps.
//
// Generated syntax:
// - Struct fields identified by a dotted path name (Struct.field)
// - Map keys identified by square brackets (Map[key])
//
// Pointers are transparently dereferenced and are not represented in the generated field name.
func listFields(obj interface{}) (fields []string) {
	return listReflectedFields(reflect.TypeOf(obj), reflect.ValueOf(obj), "")
}

// listReflectedFields recursively traverses the structure defined by reflect.Type and reflect.Value
// to discover field paths. It supports struct fields, pointer fields, and map fields (with keys).
// Fields are returned as their dotted paths. For map fields, keys are included as "[key]" segments.
//
// See listFields and getNestedFieldValue for more information.
func listReflectedFields(t reflect.Type, v reflect.Value, prefix string) (fields []string) {
	for _, field := range reflect.VisibleFields(t) {
		if !field.IsExported() || field.Anonymous {
			continue
		}
		fieldType := field.Type
		fieldValue := reflect.Value{}
		if v.IsValid() {
			fieldValue = v.FieldByName(field.Name)
		}
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
			fieldValue = fieldValue.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			fields = append(fields, listReflectedFields(fieldType, fieldValue, fmt.Sprintf("%v%v.", prefix, field.Name))...)
		} else if fieldType.Kind() == reflect.Map {
			if !fieldValue.IsValid() {
				continue
			}
			keys := fieldValue.MapKeys()
			sortValues(keys)
			for _, key := range keys {
				elementType := fieldType.Elem()
				elementValue := fieldValue.MapIndex(key)
				if elementType.Kind() == reflect.Pointer {
					elementType = elementType.Elem()
					if elementValue.IsValid() {
						elementValue = elementValue.Elem()
					}
				}
				if elementType.Kind() == reflect.Struct {
					fields = append(fields, listReflectedFields(elementType, elementValue, fmt.Sprintf("%v%v[%v].", prefix, field.Name, key.String()))...)
				} else {
					fields = append(fields, fmt.Sprintf("%v%v[%v]", prefix, field.Name, key.String()))
				}
			}
		} else {
			fields = append(fields, prefix+field.Name)
		}
	}
	return
}

// valueStr converts a reflect.Value into a string. For slices, it joins all elements into a comma-separated
// string. For pointers, it returns an empty string if the pointer is nil. For other kinds, it returns the
// formatted string representation of the value.
func valueStr(value reflect.Value) (output string) {
	if !value.IsValid() {
		return ""
	}
	if value.Kind() == reflect.Pointer {
		if value.IsZero() {
			return ""
		}
		value = value.Elem()
	}
	if value.Type() == reflect.TypeOf(net.IP{}) {
		if !value.IsZero() {
			output = fmt.Sprintf("%s", value)
		}
	} else if value.Kind() == reflect.Slice {
		var sliceStrs []string
		for i := 0; i < value.Len(); i++ {
			sliceStrs = append(sliceStrs, fmt.Sprintf("%v", value.Index(i)))
		}
		output = strings.Join(sliceStrs, ",")
	} else {
		output = fmt.Sprintf("%s", value)
	}
	return output
}

// sortValues sorts a slice of reflect.Values. Currently, it only supports sorting string values and
// will panic if values of any other kind are encountered. Values of different kinds also cannot be sorted.
func sortValues(values []reflect.Value) {
	sort.Slice(values, func(i, j int) bool {
		a, b := values[i], values[j]
		if a.Kind() != b.Kind() {
			panic(fmt.Sprintf("cannot sort values of different kinds: %s, %s", a.Kind(), b.Kind()))
		}
		switch a.Kind() {
		case reflect.String:
			return a.String() < b.String()
		default:
			panic(fmt.Sprintf("unsupported kind: %s", a.Kind()))
		}
	})
}
