package tin

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

// return the (first) type info with a name matching the passed name.
func FindByName(tins []*TypeInfo, name string) (ret *TypeInfo, okay bool) {
	for _, t := range tins {
		if t.Name == name {
			ret, okay = t, true
			break
		}
	}
	return
}

func LimitsOfNext(tins []*TypeInfo, name string) (ret block.Limits, err error) {
	if t, e := findMatchingType(tins, name, MidBlock); e != nil {
		err = e
	} else {
		ret = t.LimitsOfNext(tins)
	}
	return
}

func LimitsOfOutput(tins []*TypeInfo, name string) (ret block.Limits, err error) {
	if t, e := findMatchingType(tins, name, TermBlock); e != nil {
		err = e
	} else {
		ret, err = t.LimitsOfOutput(tins)
	}
	return
}

// find the limits of the pointer's next link
func LimitsOfType(tins []*TypeInfo, ptrType r.Type, model Model) block.Limits {
	return limitsOf(ptrType, &typeFilter{tins, model})
}

func findMatchingType(tins []*TypeInfo, typeName string, model Model) (ret *TypeInfo, err error) {
	if t, ok := FindByName(tins, typeName); !ok {
		err = errutil.New("couldnt find type '"+typeName+"' for", model)
	} else if model != t.Model {
		err = errutil.New("type '"+typeName+"' isnt", model)
	} else {
		ret = t
	}
	return
}
