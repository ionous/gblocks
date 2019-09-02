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
	wsId, blockId := main.BlockWorkspace().WorkspaceId(), main.BlockId()
	uid := block.Scope(wsId, blockId)
	ret, ok := mbs.blocks[uid]
	if !ok {
		ret = &MutatedBlock{
			block:    main,
			arch:     arch,
			atomizer: atomizer,
			inputs:   make(map[string]*MutatedInput),
			store:    make(map[string]Store),
		}
		mbs.blocks[uid] = ret
	}
	return
}

func (mbs *MutatedBlocks) GetMutatedBlock(main block.Shape) (*MutatedBlock, bool) {
	wsId, blockId := main.BlockWorkspace().WorkspaceId(), main.BlockId()
	uid := block.Scope(wsId, blockId)
	ret, ok := mbs.blocks[uid]
	return ret, ok
}

func (mbs *MutatedBlocks) OnDelete(wsId, blockId string) {
	uid := block.Scope(wsId, blockId)
	delete(mbs.blocks, uid)
}
