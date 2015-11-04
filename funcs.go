//file is borrowed from github.com/kelseyhightower/confd

package main

import (
	"encoding/json"

	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/structs"
)

func newFuncMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["base"] = path.Base
	m["split"] = strings.Split
	m["json"] = UnmarshalJsonObject
	m["jsonArray"] = UnmarshalJsonArray
	m["dir"] = path.Dir
	m["get"] = get
	m["getenv"] = os.Getenv
	m["join"] = strings.Join
	m["atoi"] = strconv.Atoi
	m["where"] = where
	m["inspect"] = Inspect
	m["datetime"] = time.Now
	m["toUpper"] = strings.ToUpper
	m["toLower"] = strings.ToLower
	m["contains"] = strings.Contains
	m["replace"] = strings.Replace
	m["scratchGetMVal"] = scratchGetMapValue
	m["scratchGetSVal"] = scratchGetSliceValue
	m["scratchSet"] = scratchSet
	m["scratchAdd"] = scratchAdd
	m["scratchMNames"] = scratchMapNames
	m["scratchSNames"] = scratchSliceNames

	return m
}

////////////////////////////////////////////////////////////////////////////////
func Inspect(args ...interface{}) {
	spew.Dump(args)
}

////////////////////////////////////////////////////////////////////////////////
func get(ctx interface{}, action string, args ...interface{}) (interface{}, error) {
	method := reflect.ValueOf(ctx).MethodByName(action)
	in := make([]reflect.Value, len(args))
	for idx, arg := range args {
		in[idx] = reflect.ValueOf(arg)
	}

	out := method.Call(in)
	ret := out[0].Interface()
	err := out[1].Interface()

	if err != nil {
		return ret, err.(error)
	}
	return ret, nil
}

////////////////////////////////////////////////////////////////////////////////
func where(in interface{}, field string, sliceVal interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	if in == nil {
		return ret, errors.New("where: source is nil")
	}
	if field == "" {
		return ret, errors.New("where: field is empty")
	}
	if sliceVal == nil {
		return ret, errors.New("where: value is nil")
	}

	if reflect.TypeOf(in).Kind() != reflect.Slice {
		return ret, errors.New("where: source is no slice value")
	}

	s := reflect.ValueOf(in)
	for i := 0; i < s.Len(); i++ {
		val := s.Index(i).Interface()
		st := structs.New(val)
		fieldVal, ok := st.FieldOk(field)
		if !ok {
			return ret, errors.Errorf("where: key %q not found", field)
		}
		if fieldVal.Value() == sliceVal {
			ret = append(ret, val)
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
