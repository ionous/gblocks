package mutant

import "strconv"

var _autoid int

func TempNewId() string {
	ret := "idgen" + strconv.Itoa(_autoid)
	_autoid++
	return ret
}
