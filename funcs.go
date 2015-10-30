//file is borrowed from github.com/kelseyhightower/confd

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
)

func newFuncMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["base"] = path.Base
	m["split"] = strings.Split
	m["json"] = UnmarshalJsonObject
	m["jsonArray"] = UnmarshalJsonArray
	m["dir"] = path.Dir
	m["getenv"] = os.Getenv
	m["join"] = strings.Join
	m["atoi"] = strconv.Atoi
	m["where"] = where
	m["datetime"] = time.Now
	m["toUpper"] = strings.ToUpper
	m["toLower"] = strings.ToLower
	m["contains"] = strings.Contains
	m["replace"] = strings.Replace
	return m
}

// ToSliceE casts an empty interface to a []interface{}.
func ToSliceE(i interface{}) ([]interface{}, error) {
	printDebug("ToSliceE called on type:", reflect.TypeOf(i))

	var s []interface{}

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	case []map[string]interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	default:
		return s, fmt.Errorf("Unable to Cast %#v of type %v to []interface{}", i, reflect.TypeOf(i))
	}
}

////////////////////////////////////////////////////////////////////////////////
func where(in interface{}, sliceKey string, sliceVal interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	if in == nil {
		return ret, errors.New("where: source is nil")
	}
	if sliceKey == "" {
		return ret, errors.New("where: key is nil")
	}
	if sliceVal == nil {
		return ret, errors.New("where: value is nil")
	}

	m, err := ToSliceE(in)
	if err != nil {
		return ret, err
	}

	for _, str := range m {
		st := structs.New(str)
		field, ok := st.FieldOk(sliceKey)
		if !ok {
			return ret, errors.New("where: input is no []interface{} value")
		}
		if field.Value() == sliceVal {
			ret = append(ret, str)
		}
	}

	return ret, nil
}

func addFuncs(out, in map[string]interface{}) {
	for name, fn := range in {
		out[name] = fn
	}
}

func UnmarshalJsonObject(data string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func UnmarshalJsonArray(data string) ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}
