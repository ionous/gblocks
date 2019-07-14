package mutant

import "github.com/ionous/gblocks/block"

// create a new of mutable blocks.
// these blocks are editable by the user using the mutation ui.
func NewMutatedBlocks() *MutatedBlocks {
	return &MutatedBlocks{make(map[string]*mutatedBlock)}
}

// maps workspace blocks to user edited mutation data.
// shared across every Mutator.
type MutatedBlocks struct {
	blocks map[string]*mutatedBlock
}

func (mbs *MutatedBlocks) AddMutatedBlock(main block.Shape, atoms AtomizedInputs) {
	wid := main.BlockId()
	mbs.blocks[wid] = &mutatedBlock{main, atoms, nil}
}

func (mbs *MutatedBlocks) GetMutatedBlock(main block.Shape) (*mutatedBlock, bool) {
	wid := main.BlockId()
	ret, ok := mbs.blocks[wid]
	return ret, ok
}

func (mbs *MutatedBlocks) EnsureMutatedBlock(main block.Shape) *mutatedBlock {
	wid := main.BlockId()
	ret, ok := mbs.blocks[wid]
	if !ok {
		ret = &mutatedBlock{block: main}
		mbs.blocks[wid] = ret
	}
	return ret
}

func (mbs *MutatedBlocks) OnDelete(wid string) {
	delete(mbs.blocks, wid)
}
