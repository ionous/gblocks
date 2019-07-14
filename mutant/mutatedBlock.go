package mutant

import "github.com/ionous/gblocks/block"

// each mutable workspace block contains
// one or more mutable inputs containing a list of atoms
// and connections to existing workspace blocks.
type mutatedBlock struct {
	block  block.Shape // workspace block
	atoms  AtomizedInputs
	stored SavedConnections
}

func (mb *mutatedBlock) ContainerName() string {
	return ContainerName(mb.block.BlockType())
}

func (mb *mutatedBlock) GetAtomsForInput(inputName string) ([]string, bool) {
	return mb.atoms.GetAtomsForInput(inputName)
}

func (mb *mutatedBlock) RemoveAtoms() {
	RemoveAtoms(mb.block)
}

func (mb *mutatedBlock) SaveConnections(mui block.Shape) {
	mb.stored = SaveConnections(mb.block, mui)
}

func (mb *mutatedBlock) RestoreConnections() {
	mb.stored.RestoreConnections(mb.block)
}
