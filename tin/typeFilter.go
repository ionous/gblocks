package tin

import (
	r "reflect"
)

// iterator over a slice of type infos
type typeFilter struct {
	tins   []*TypeInfo // slice collapses as we go
	filter Model       // only produce tins which match this model
}

func (it *typeFilter) Name() string {
	t := it.TypeInfo()
	return t.Name
}

func (it *typeFilter) PtrType() r.Type {
	t := it.TypeInfo()
	return t.ptrType
}

func (it *typeFilter) TypeInfo() *TypeInfo {
	return it.tins[0]
}

func (it *typeFilter) NextType() (ret typeIterator, okay bool) {
	for i, cnt := 1, len(it.tins); i < cnt; i++ {
		if t := it.tins[i]; t.Model == it.filter {
			ret, okay = &typeFilter{it.tins[i:], it.filter}, true
			break
		}
	}
	return
}
