package toolbox

import (
	r "reflect"
	"strconv"
)

// see if the type implements the stringer, ex. an enum.
func asStringer(t r.Type, v r.Value) (ret string, okay bool) {
	type stringer interface{ String() string }
	stringish := t.AssignableTo(r.TypeOf((*stringer)(nil)).Elem())
	if stringish {
		if !v.IsValid() {
			v = r.Zero(t)
		}
		if val, ok := v.Interface().(stringer); ok {
			ret = val.String()
		}
		okay = true // supports stringer, even if nil
	}

	return
}

func asString(v r.Value) (ret string, okay bool) {
	if v.IsValid() {
		if val := v.String(); len(val) > 0 {
			ret, okay = val, true
		}
	}
	return
}

func asBool(v r.Value) (ret string, okay bool) {
	if v.IsValid() {
		if val := v.Bool(); val {
			ret = strconv.FormatBool(val)
			okay = true
		}
	}
	return
}

func asInt(v r.Value) (ret string, okay bool) {
	if v.IsValid() {
		if val := v.Int(); val != 0 {
			ret = strconv.FormatInt(val, 10)
			okay = true
		}
	}
	return
}

func asUint(v r.Value) (ret string, okay bool) {
	if v.IsValid() {
		if val := v.Uint(); val != 0 {
			ret = strconv.FormatUint(val, 10)
			okay = true
		}
	}
	return
}

func asFloat(v r.Value) (ret string, okay bool) {
	if v.IsValid() {
		if val := v.Float(); val != 0 {
			ret = strconv.FormatFloat(val, 'g', -1, 32)
			okay = true
		}
	}
	return
}
