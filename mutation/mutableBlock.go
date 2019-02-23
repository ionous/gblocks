package mutation

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/blockly"
)

// MutationType -  name of the container used for the mutator ui.
func MutationType(typeName block.Type) block.Type {
	return block.SpecialType("mui_container", typeName.String())
}

// associates mutableBlocks with workspace and block pointers; alleviates passing a large numbers of parameters to helper functions
type mutableBlock struct {
	workspace *WorkspaceFactory
	block     *blockly.Block
	blockData *blockData
}

func (mb *mutableBlock) mutationByIndex(inputIndex int) (ret *mutableInput, okay bool) {
	if inputIndex < mb.block.NumInputs() {
		in := mb.block.Input(inputIndex)
		if md, ok := mb.blockData.inputMap[in.Name]; ok {
			ret, okay = &mutableInput{mb, in.Name, inputIndex, md}, true
		}
	}
	return
}

func (mb *mutableBlock) mutationByName(name block.Item) (ret *mutableInput, index int) {
	index = -1 // not found
	if in, inputIndex := mb.block.InputByName(name); inputIndex >= 0 {
		if md, ok := mb.blockData.inputMap[in.Name]; ok {
			ret, index = &mutableInput{mb, in.Name, inputIndex, md}, inputIndex
		}
	}
	return
}

// collapse all dynamic inputs
func (mb *mutableBlock) removeAtoms() {
	var collapse int
	block := mb.block
	for i := 0; i < block.NumInputs(); {
		in := block.Input(i)
		if collapse > 0 {
			block.RemoveInput(in.Name)
			collapse--
		} else {
			if mi, ok := mb.mutationByIndex(i); ok {
				collapse = mi.resetAtoms() // returns the total number of inputs used by this mutation
			}
			i++ // dont advance i if we collapse.
		}
	}
}
