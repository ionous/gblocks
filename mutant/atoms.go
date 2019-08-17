package mutant

import (
	"strings"

	"github.com/ionous/gblocks/block"
)

func RemoveAtoms(b block.Inputs) {
	atomPrefix := block.Scope("a", "")
	for inputIndex := 0; inputIndex < b.NumInputs(); {
		in := b.Input(inputIndex)
		inputName := in.InputName()
		if strings.HasPrefix(inputName, atomPrefix) {
			b.RemoveInput(inputName)
		} else {
			inputIndex++
		}
	}
}
