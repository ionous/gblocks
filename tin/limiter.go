package tin

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

// iterate over a list of block types to determine limits, for example:
// . workspace block stacking ( iterating over statement tins )
// . workspace links between between inputs and outputs ( iterating over term tins )
// . mui blocks (quarks) stacking.
type typeIterator interface {
	Name() string
	PtrType() r.Type
	NextType() (typeIterator, bool)
}

// find the limits of the pointer's next link
func limitsOfNext(ptrType r.Type, it typeIterator) (ret block.Limits) {
	elm := ptrType.Elem()
	if f, ok := elm.FieldByName(block.NextStatement); ok {
		ret = limitsOf(f.Type, it)
	}
	return
}

// find the limits of the pointer's next link
func limitsOfOutput(ptrType r.Type, it typeIterator) (ret block.Limits, err error) {
	if outType, e := _outputType(ptrType); e != nil {
		err = e
	} else if outType != nil {
		if l := limitsOf(outType, it); !l.Connects {
			err = errutil.New("couldnt determine constraints for output", outType)
		} else {
			ret = l
		}
	}
	return
}

// what can attach to the passed slot
func limitsOf(slotType r.Type, it typeIterator) (ret block.Limits) {
	if slotType != nil {
		switch k := slotType.Kind(); k {
		case r.Ptr:
			if elType := slotType.Elem(); elType.Kind() == r.Struct {
				ret = _limitsOf(slotType, it)
			}
		case r.Interface:
			// optimization to check for interface{} -- no limits
			if basicInterface := r.TypeOf((*interface{})(nil)).Elem(); slotType == basicInterface {
				ret = block.MakeUnlimited()
			} else {
				ret = _limitsOf(slotType, it)
			}
		}
	}
	return
}

func _limitsOf(slotType r.Type, it typeIterator) (ret block.Limits) {
	var unassigned int
	var types []string
	for okay := true; okay; it, okay = it.NextType() {
		// check that the type
		if ptrType := it.PtrType(); ptrType.Kind() == r.Ptr {
			if !ptrType.AssignableTo(slotType) {
				unassigned++
			} else {
				// note: our types can include interface, but interfaces are not blocks
				// and cannot actually be assigned ( even if they are in fact compatible
				types = append(types, it.Name())
			}
		}
	}
	if unassigned == 0 {
		// if nothing was unassigned; everything was assigned; we can return unliited
		ret = block.MakeUnlimited()
	} else if len(types) > 0 {
		ret = block.MakeLimits(types)
	} else {
		ret = block.MakeOffLimits()
	}
	return
}

// if there is no output connection return is nil
func _outputType(ptrType r.Type) (ret r.Type, err error) {
	if output, ok := ptrType.MethodByName("Output"); ok {
		if cnt := output.Type.NumOut(); cnt != 1 {
			err = errutil.New("unexpected output count", cnt)
		} else if t := output.Type.Out(0); Classify(t) != InputClass {
			err = errutil.New("unexpected output type", t)
		} else {
			ret = t
		}
	}
	return
}
