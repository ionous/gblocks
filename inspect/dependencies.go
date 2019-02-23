package inspect

import (
	"github.com/ionous/gblocks/block"
	r "reflect"
)

type DependencyPool struct {
	types map[block.Type]r.Type //  *Struct
}

func (dp *DependencyPool) AddType(ptrType r.Type) {
	if ptrType.Kind() != r.Ptr && ptrType.Elem().Kind() != r.Struct {
		panic("expected ptr to struct; was: " + ptrType.String())
	}
	name := block.TypeFromStruct(ptrType.Elem())
	if dp.types == nil {
		dp.types = map[block.Type]r.Type{name: ptrType}
	} else {
		dp.types[name] = ptrType
	}
}

func (dp *DependencyPool) GetConstraints(slotType r.Type) ([]block.Type, bool) {
	var constraints []block.Type
	var hasConnection bool
	if basicInterface := r.TypeOf((interface{})(nil)); slotType == basicInterface {
		hasConnection = true
	} else {
		for name, ptrType := range dp.types {
			if ptrType.AssignableTo(slotType) {
				constraints = append(constraints, name)
				hasConnection = true
			}
		}
	}
	return constraints, hasConnection
}
