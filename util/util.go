package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

// MapToStruct is a way to deserialize a dictionary of strings into a struct
//
func MapHeaderToStruct(m map[string][]string, val interface{}) error {
	simple := make(map[string]string, len(m))
	for k, v := range m {
		simple[k] = v[0]
	}
	return MapToStruct(simple, &val)
}

// MapToStruct is a way to deserialize a dictionary of strings into a struct
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

func NewUid() string {
	uid := uuid.NewRandom()
	return Base64Safe(uid)
}

func NoQuotes(target string) string {
	return strings.Replace(target, "\"", "", -1)
}

// Base64Safe - Makes safe for file paths, removes the forward slash '/' and equals '='
func Base64Safe(bts []byte) string {
	safe := strings.Replace(base64.StdEncoding.EncodeToString(bts), "/", "-", -1)
	return strings.Replace(safe, "=", "", -1)
}

func NotEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func NotNil(values ...interface{}) interface{} {
	for _, v := range values {
		if v != nil {
			return v
		}
	}
	return nil

}
