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
)

func newFuncMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["base"] = path.Base
	m["split"] = strings.Split
	m["json"] = UnmarshalJsonObject
	m["jsonArray"] = UnmarshalJsonArray
	m["dir"] = path.Dir
	m["getenv"] = os.Getenv
	m["sliceselect"] = sliceselect
	m["join"] = strings.Join
	m["atoi"] = strconv.Atoi
	m["datetime"] = time.Now
	m["toUpper"] = strings.ToUpper
	m["toLower"] = strings.ToLower
	m["contains"] = strings.Contains
	m["replace"] = strings.Replace
	return m
}

////////////////////////////////////////////////////////////////////////////////
func sliceselect(m []interface{}, sliceKey string, sliceVal interface{}) ([]interface{}, error) {
	if m == nil {
		ret := make([]interface{}, 0)
		return ret, errors.New("sliceselect: first argument is nil")
	}
	if sliceKey == "" {
		ret := make([]interface{}, 0)
		return ret, errors.New("sliceselect: second argument is nil")
	}
	if sliceVal == nil {
		ret := make([]interface{}, 0)
		return ret, errors.New("sliceselect: third argument is nil")
	}

	idx := 0
	v := make([]map[string]interface{}, len(m))
	for _, val := range m {
		v[idx] = val.(map[string]interface{})
		idx++
	}

	v2 := make([]interface{}, 0)
	for _, mp := range v {
		if s, ok := mp[sliceKey]; ok {
			for _, vl := range s.([]interface{}) {
				if vl == sliceVal {
					v2 = append(v2, mp)
				}
			}
		}
	}

	return v2, nil
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
