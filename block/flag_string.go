// Code generated by "stringer -type=Flag"; DO NOT EDIT.

package block

import "fmt"

const _Flag_name = "DeletableMovableShadowInsertionMakerEditableCollapsedEnabledInputsInline"

var _Flag_index = [...]uint8{0, 9, 16, 22, 36, 44, 53, 60, 72}

func (i Flag) String() string {
	if i < 0 || i >= Flag(len(_Flag_index)-1) {
		return fmt.Sprintf("Flag(%d)", i)
	}
	return _Flag_name[_Flag_index[i]:_Flag_index[i+1]]
}