package main

var scrtsMap = make(map[string]map[string]interface{})
var scrtsSlice = make(map[string][]interface{})

//////////////////////////////////////////////////////////////////////////////
func scratchMapSet(scratch string, key string, value interface{}) string {
	if scr, ok := scrtsMap[scratch]; ok {
		scr[key] = value
	} else {
		scrtsMap[scratch] = make(map[string]interface{})
		scrtsMap[scratch][key] = value
	}

	return ""
}

//////////////////////////////////////////////////////////////////////////////
func scratchSliceAdd(scratch string, value interface{}) string {
	if scr, ok := scrtsSlice[scratch]; ok {
		scr = append(scr, value)
	} else {
		scrtsSlice[scratch] = make([]interface{}, 0)
		scrtsSlice[scratch] = append(scrtsSlice[scratch], value)
	}

	return ""
}

//////////////////////////////////////////////////////////////////////////////
func scratchGetMapValue(scratch string, key string) interface{} {
	if scr, ok := scrtsMap[scratch]; ok {
		return scr[key]
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////
func scratchGetSliceValues(scratch string) []interface{} {
	if scr, ok := scrtsSlice[scratch]; ok {
		return scr
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////
func scratchMapNames() []string {
	var nms []string
	for k := range scrtsMap {
		nms = append(nms, k)
	}
	return nms
}

//////////////////////////////////////////////////////////////////////////////
func scratchSliceNames() []string {
	var nms []string
	for k := range scrtsSlice {
		nms = append(nms, k)
	}
	return nms
}
