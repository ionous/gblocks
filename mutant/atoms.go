package mutant

import (
	"strings"

	"github.com/ionous/gblocks/block"
)

// MutableInputs --  per mutable input, a listing the user's selected atoms
type MutableInputs map[string][]string


func RemoveAtoms(b block.Inputs) {
	prefix := block.Scope("a", "")
	for inputIndex := 0; inputIndex < b.NumInputs(); {
		in := b.Input(inputIndex)
		inputName := in.InputName()
		if strings.HasPrefix(inputName, prefix) {
			b.RemoveInput(inputName)
		} else {
			inputIndex++
		}
	}
}
