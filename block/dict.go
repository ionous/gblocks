package block

type Dict map[string]interface{}

func (d Dict) Lookup(key string) (ret string, okay bool) {
	if val, ok := d[key]; ok {
		ret, okay = val.(string)
	}
	return
}

// add key to this from dict or the defaultValue if it doesnt exist
func Merge(dst, src Dict, key string, defaultValue interface{}) {
	if val, ok := src[key]; ok {
		dst[key] = val
	} else {
		dst[key] = defaultValue
	}
}
