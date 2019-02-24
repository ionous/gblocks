package inspect

import (
	"github.com/ionous/errutil"
	r "reflect"
	"sort"
)

type EnumPair []string // display, uniquifier

type EnumPairs struct {
	mapping map[r.Type][]EnumPair
}
type stringer interface{ String() string }

// returns nil if nothing registered
func (ep EnumPairs) GetPairs(rtype r.Type) []EnumPair {
	return ep.mapping[rtype]
}

// a map of enum value to string
func (ep *EnumPairs) AddEnum(mapping interface{}) (ret []EnumPair, err error) {
	if src, srcType := r.ValueOf(mapping), r.TypeOf(mapping); srcType.Kind() != r.Map {
		err = errutil.New("invalid enum mapping", srcType)
	} else {
		keyType, valueType := srcType.Key(), srcType.Elem()
		if classify(keyType) != Int {
			err = errutil.New("invalid enum key type", keyType)
		} else if valueType.Kind() != r.String {
			err = errutil.New("invalid enum value type", valueType)
		} else {
			pairs := makePairs(src)
			if ep.mapping == nil {
				ep.mapping = map[r.Type][]EnumPair{keyType: pairs}
			} else {
				ep.mapping[keyType] = pairs
			}
			ret = pairs
		}
	}
	return
}

func makePairs(src r.Value) (ret []EnumPair) {
	if keys := src.MapKeys(); len(keys) > 0 {
		var pairs []EnumPair
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
			pair := EnumPair{display, unique}
			pairs = append(pairs, pair)
		}
		ret = pairs
	}
	return
}
