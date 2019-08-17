package mutant

import "github.com/ionous/gblocks/block"

// create a new of mutable blocks.
// these blocks are editable by the user using the mutation ui.
func NewMutatedBlocks() *MutatedBlocks {
	return &MutatedBlocks{make(map[string]*MutatedBlock)}
}

// maps workspace blocks to user edited mutation data.
// shared across every Mutator.
type MutatedBlocks struct {
	blocks map[string]*MutatedBlock
}

func (mbs *MutatedBlocks) CreateMutatedBlock(main block.Shape, arch *BlockMutations, atomizer Atomizer) (ret *MutatedBlock) {
	wid := main.BlockId()
	ret, ok := mbs.blocks[wid]
	if !ok {
		ret = &MutatedBlock{
			block:    main,
			arch:     arch,
			atomizer: atomizer,
			inputs:   make(map[string]*MutatedInput),
			store:    make(map[string]Store),
		}
		mbs.blocks[wid] = ret
	}
	return
}

func (mbs *MutatedBlocks) GetMutatedBlock(main block.Shape) (*MutatedBlock, bool) {
	wid := main.BlockId()
	ret, ok := mbs.blocks[wid]
	return ret, ok
}

func (mbs *MutatedBlocks) OnDelete(wid string) {
	delete(mbs.blocks, wid)
}
