package mutant

import (
	"strings"

	"github.com/ionous/gblocks/block"
)

// AtomizedInputs --  per mutable input, a listing the user's selected atoms
type AtomizedInputs map[string][]string

func MakeAtomizedInputs() AtomizedInputs {
	return make(AtomizedInputs)
}

func (mi *AtomizedInputs) GetAtomsForInput(inputName string) ([]string, bool) {
	ret, ok := (*mi)[inputName]
	return ret, ok
}

func (mi *AtomizedInputs) SetAtomsForInput(inputName string, atoms []string) {
	(*mi)[inputName] = atoms
}

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
