package gblocks

import (
	"github.com/ionous/errutil"
	r "reflect"
	"strings"
)

// RegisteredTypes -- IsRegistered structs.
type RegisteredTypes map[TypeName]r.Type

// IsRegistered -
func (ts RegisteredTypes) IsRegistered(typeName TypeName) (okay bool) {
	if _, exists := ts[typeName]; exists {
		okay = true
	}
	return
}

// RegisterType - blockType should generally be a struct
func (ts RegisteredTypes) RegisterType(blockType r.Type) (newlyAdded bool) {
	if typeName := toTypeName(blockType); !ts.IsRegistered(typeName) {
		ts[typeName] = blockType
		newlyAdded = true
	}
	return
}

// CheckField - return the types which can satisf the passed field.
func (ts RegisteredTypes) CheckField(structType r.Type, field string) (ret Constraints, err error) {
	if f, ok := structType.FieldByName(field); !ok {
		// no error, zero value of constraints means no contrains
	} else {
		ret, err = ts.CheckStructField(f)
	}
	return
}

func (ts RegisteredTypes) CheckStructField(f r.StructField) (ret Constraints, err error) {
	if tag, ok := f.Tag.Lookup(opt_check); !ok {
		ret, err = ts.CheckType(f.Type)
	} else {
		// read the tag string
		if tagParts := strings.Split(tag, ","); len(tagParts) > 0 {
			for _, n := range tagParts {
				n := pascalToUnderscore(strings.TrimSpace(n))
				if typeName := TypeName(n); !ts.IsRegistered(typeName) {
					err = errutil.New("unknown type in constraint", typeName)
				} else {
					ret.AddConstraint(typeName)
				}
			}
		}
	}
	return
}

// CheckType - return the types which can satisfy type t.
func (ts RegisteredTypes) CheckType(t r.Type) (ret Constraints, err error) {
	switch t.Kind() {
	case r.Ptr:
		if elType := t.Elem(); elType.Kind() != r.Struct {
			err = errutil.New("unexpected type", elType)
		} else if typeName := toTypeName(elType); !ts.IsRegistered(typeName) {
			err = errutil.New("unknown pointer type", t)
		} else {
			ret.AddConstraint(typeName)
		}
	case r.Interface:
		if basicInterface := r.TypeOf((interface{})(nil)); t == basicInterface {
			ret.AddConstraint("")
		} else {
			for typeName, srcType := range ts {
				// registered types hold structType, for implementation we expect func (*struct) interface{}.
				if r.PtrTo(srcType).Implements(t) {
					ret.AddConstraint(typeName)
				}
			}
		}
	default:
		err = errutil.New("unknown connection type", t)
	}
	return
}
