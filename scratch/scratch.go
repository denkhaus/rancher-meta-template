package scratch

var (
	scrtsMap   = make(map[string]map[string]interface{})
	scrtsSlice = make(map[string][]interface{})
)

//////////////////////////////////////////////////////////////////////////////
func Reset() {
	scrtsMap = make(map[string]map[string]interface{})
	scrtsSlice = make(map[string][]interface{})
}

//////////////////////////////////////////////////////////////////////////////
func ScratchMapSet(scratch string, key string, value interface{}) string {
	if scr, ok := scrtsMap[scratch]; ok {
		scr[key] = value
		scrtsMap[scratch] = scr
	} else {
		scrtsMap[scratch] = make(map[string]interface{})
		scrtsMap[scratch][key] = value
	}

	return ""
}

//////////////////////////////////////////////////////////////////////////////
func ScratchSliceAdd(scratch string, value interface{}) string {
	if scr, ok := scrtsSlice[scratch]; ok {
		scr = append(scr, value)
		scrtsSlice[scratch] = scr
	} else {
		scrtsSlice[scratch] = []interface{}{value}
	}

	return ""
}

//////////////////////////////////////////////////////////////////////////////
func ScratchGetMapValue(scratch string, key string) interface{} {
	if scr, ok := scrtsMap[scratch]; ok {
		return scr[key]
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////
func ScratchGetSliceValues(scratch string) []interface{} {
	if scr, ok := scrtsSlice[scratch]; ok {
		return scr
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////
func ScratchMapNames() []string {
	var nms []string
	for k := range scrtsMap {
		nms = append(nms, k)
	}
	return nms
}

//////////////////////////////////////////////////////////////////////////////
func ScratchSliceNames() []string {
	var nms []string
	for k := range scrtsSlice {
		nms = append(nms, k)
	}
	return nms
}
