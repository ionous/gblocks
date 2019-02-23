package block

type Dict map[string]interface{}

func (d Dict) Contains(key string) (okay bool) {
	if _, ok := d[key]; ok {
		okay = true
	}
	return
}

// add only if the key is new
func (d Dict) Insert(key string, value interface{}) {
	if !d.Contains(key) {
		d[key] = value
	}
}

// shallow
func (d Dict) Copy() Dict {
	out := make(Dict)
	for k, v := range d {
		out[k] = v
	}
	return out
}
