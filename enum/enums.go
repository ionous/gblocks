package enum

import (
	r "reflect" // for inspecting enumerated maps
	"sort"

	"github.com/ionous/errutil"
)

type Pair []string // display, uniquifier

type Pairs struct {
	mapping map[string][]Pair
}
type stringer interface{ String() string }

// name is generally pascal cased; returns nil if nothing registered
func (ep *Pairs) GetPairs(name string) []Pair {
	return ep.mapping[name]
}

// a map of enum value to string
func (ep *Pairs) AddEnum(mapping interface{}) (ret []Pair, err error) {
	if src, srcType := r.ValueOf(mapping), r.TypeOf(mapping); srcType.Kind() != r.Map {
		err = errutil.New("invalid enum mapping", srcType)
	} else {
		keyType, valueType := srcType.Key(), srcType.Elem()
		if keyType.Kind() != r.Int {
			err = errutil.New("invalid enum key type", keyType)
		} else if valueType.Kind() != r.String {
			err = errutil.New("invalid enum value type", valueType)
		} else {
			pairs := makePairs(src)
			name := keyType.Name()
			if ep.mapping == nil {
				ep.mapping = map[string][]Pair{name: pairs}
			} else {
				ep.mapping[name] = pairs
			}
			ret = pairs
		}
	}
	return
}

func makePairs(src r.Value) (ret []Pair) {
	if keys := src.MapKeys(); len(keys) > 0 {
		var pairs []Pair
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].Int() < keys[j].Int()
		})
		for _, key := range keys {
			var unique string
			if stringer, ok := key.Interface().(stringer); ok {
				//unique := fmt.Sprint(key)
				unique = stringer.String()
			} else {
				unique = key.String()
			}
			display := src.MapIndex(key).String()
			pair := Pair{display, unique}
			pairs = append(pairs, pair)
		}
		ret = pairs
	}
	return
}
