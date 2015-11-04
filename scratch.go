package main

var scrtsMap = make(map[string]map[string]interface{})
var scrtsSlice = make(map[string][]interface{})

//////////////////////////////////////////////////////////////////////////////
func scratchSet(scratch string, key string, value interface{}) {
	if scr, ok := scrtsMap[scratch]; ok {
		scr[key] = value
	} else {
		scrtsMap[scratch] = make(map[string]interface{})
		scrtsMap[scratch][key] = value
	}
}

//////////////////////////////////////////////////////////////////////////////
func scratchAdd(scratch string, value interface{}) {
	if scr, ok := scrtsSlice[scratch]; ok {
		scr = append(scr, value)
	} else {
		scrtsSlice[scratch] = make([]interface{}, 0)
		scrtsSlice[scratch] = append(scrtsSlice[scratch], value)
	}
}

//////////////////////////////////////////////////////////////////////////////
func scratchGetMapValue(scratch string, key string) interface{} {
	if scr, ok := scrtsMap[scratch]; ok {
		return scr[key]
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////
func scratchGetSliceValue(scratch string, key string) []interface{} {
	if scr, ok := scrtsSlice[scratch]; ok {
		return scr
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////
func scratchMapNames() []string {
	nms := make([]string, 0)

	for k, _ := range scrtsMap {
		nms = append(nms, k)
	}
	return nms
}

//////////////////////////////////////////////////////////////////////////////
func scratchSliceNames() []string {
	nms := make([]string, 0)

	for k, _ := range scrtsSlice {
		nms = append(nms, k)
	}
	return nms
}
