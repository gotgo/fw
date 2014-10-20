package util

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// MapToStruct is a way to deserialize a dictionary of strings into a struct
// To enable a terse syntax the following pattern is prefered:
//
//		myType, err := MapToStruct(source, new (MyType))
//
func MapHeaderToStruct(m map[string][]string, val interface{}) error {
	simple := make(map[string]string, len(m))
	for k, v := range m {
		simple[k] = v[0]
	}
	return MapToStruct(simple, &val)
}

// MapToStruct is a way to deserialize a dictionary of strings into a struct
// To enable a terse syntax the following pattern is prefered:
//
//		myType, err := MapToStruct(source, new (MyType))
//
func MapToStruct(m map[string]string, val interface{}) error {
	b := new(bytes.Buffer)
	e := json.NewEncoder(b)
	e.Encode(m)
	d := json.NewDecoder(b)
	err := d.Decode(&val)
	return err
}

//      Invoke(YourT2{}, "MethodFoo", 10, "abc")
//      Invoke(YourT1{}, "MethodBar")
func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}
