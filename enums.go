package gblocks

import (
	"github.com/ionous/errutil"
	r "reflect"
	"sort"
)

type EnumPair [2]string // display, uniquifier

type RegisteredEnum struct {
	pairs []EnumPair
}

type RegisteredEnums struct {
	typeToEnum map[TypeName]*RegisteredEnum
}

func (reg *RegisteredEnums) GetEnum(typeName TypeName) (ret *RegisteredEnum, okay bool) {
	if enumType, ok := reg.typeToEnum[typeName]; ok {
		ret = enumType
		okay = true
	}
	return
}

func (reg *RegisteredEnums) registerEnum(n interface{}) (ret []EnumPair, err error) {
	var pairs []EnumPair
	if src, srcType := r.ValueOf(n), r.TypeOf(n); srcType.Kind() != r.Map {
		err = errutil.New("invalid enum mapping")
	} else if keyType, valueType := srcType.Key(), srcType.Elem(); valueType.Kind() != r.String {
		err = errutil.New("invalid enum value type")
	} else {
		// want to build an array of display to stringer
		// want to sort that array for display
		// want to store that at the v.Type for lookup.
		// eventually a map? probably of stringer to int ( for reverse conversion, setting in response to changes )
		keys := src.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].Int() < keys[j].Int()
		})

		type stringer interface{ String() string }

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
		if reg.typeToEnum == nil {
			reg.typeToEnum = make(map[TypeName]*RegisteredEnum)
		}
		enumName := toTypeName(keyType)
		reg.typeToEnum[enumName] = &RegisteredEnum{pairs: pairs}
		ret = pairs
	}
	return
}
