//file is borrowed from github.com/kelseyhightower/confd

package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"
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

	m := in.([]interface{})

	for _, str := range m {
		st := structs.New(str)
		field := st.Field(sliceKey)
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
