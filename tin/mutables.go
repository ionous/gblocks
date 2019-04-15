package tin

import (
	r "reflect"

	"github.com/ionous/errutil"
)

type Mutables struct {
	slice []*Mutable
}

// mutant is a nil pointer to mutation struct (*BlockMutation)(nil).
// quarks lists the blocks of "atoms" which add items to a block dynamically
func (ms *Mutables) AddMutation(mutant interface{}, quarks ...interface{}) error {
	_, e := ms.addMutation(mutant, quarks...)
	return e
}

func (ms *Mutables) addMutation(mutant interface{}, quarks ...interface{}) (ret *Mutable, err error) {
	// get some temporary type info for the passed pointer
	if t, e := UnknownModel.PtrInfo(mutant); e != nil {
		err = errutil.New("error inspecting mutation", e)
	} else if was, found := ms.FindMutable(t.Name); found {
		err = errutil.New("mutation already registered", was)
	} else {
		mutation := &Mutable{t.Name, t.ptrType, make([]r.Type, len(quarks))}
		for i, q := range quarks {
			ptrType := r.TypeOf(q)
			if Classify(ptrType) != InputClass {
				e := errutil.New("unknown block pointer", q)
				err = errutil.Append(err, e)
			} else {
				mutation.quarks[i] = ptrType
			}
		}
		if err == nil {
			ms.slice = append(ms.slice, mutation)
			ret = mutation
		}
	}
	return
}

func (ms *Mutables) FindMutable(name string) (ret *Mutable, okay bool) {
	for _, m := range ms.slice {
		if m.name == name {
			ret, okay = m, true
			break
		}
	}
	return
}
