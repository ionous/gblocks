package inspect

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	r "reflect"
)

type DependencyPool struct {
	Types map[block.Type]r.Type //  *Struct
}

func (dp *DependencyPool) AddType(ptrType r.Type) (err error) {
	if ptrType.Kind() != r.Ptr && ptrType.Elem().Kind() != r.Struct {
		err = errutil.New("expected ptr to struct; was: " + ptrType.String())
	} else {
		name := block.TypeFromStruct(ptrType.Elem())
		if dp.Types == nil {
			dp.Types = map[block.Type]r.Type{name: ptrType}
		} else {
			dp.Types[name] = ptrType
		}
	}
	return
}

func (dp *DependencyPool) AddTypes(ptrTypes ...r.Type) (err error) {
	for _, ptrType := range ptrTypes {
		if e := dp.AddType(ptrType); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (dp *DependencyPool) GetConstraints(slotType r.Type) ([]block.Type, bool) {
	var constraints []block.Type
	var hasConnection bool

	if basicInterface := r.TypeOf((*interface{})(nil)).Elem(); slotType == basicInterface {
		hasConnection = true
	} else {
		for name, ptrType := range dp.Types {
			if ptrType.AssignableTo(slotType) {
				constraints = append(constraints, name)
				hasConnection = true
			}
		}
	}
	return constraints, hasConnection
}
