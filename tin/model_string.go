// Code generated by "stringer -type=Model"; DO NOT EDIT.

package tin

import "fmt"

const _Model_name = "UnknownModelTopBlockMidBlockTermBlock"

var _Model_index = [...]uint8{0, 12, 20, 28, 37}

func (i Model) String() string {
	if i < 0 || i >= Model(len(_Model_index)-1) {
		return fmt.Sprintf("Model(%d)", i)
	}
	return _Model_name[_Model_index[i]:_Model_index[i+1]]
}
