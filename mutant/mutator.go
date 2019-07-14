package mutant

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
)

// instance data for mutations
// implements block.Mutator
type Mutator struct {
	mins     *BlockMutations // descriptions of atoms
	atomizer Atomizer        // needed to expand quarks into atoms
	// everyone has the same atomizer, so consider putting in the mutated blocks
	blockPool *MutatedBlocks // per-block mutation data (pooled across mutators)
}

func NewMutator(mins *BlockMutations, db Atomizer, mbs *MutatedBlocks) *Mutator {
	return &Mutator{mins, db, mbs}
}

func (a *Mutator) Delete(id string) {
	a.blockPool.OnDelete(id)
}

func (a *Mutator) MutationToDom(main block.Shape) (ret string, err error) {
	if src, ok := a.blockPool.GetMutatedBlock(main); ok {
		dom := a.mins.SaveMutation(src.atoms)
		ret, err = dom.MarshalMutation()
		//println("mutation to dom", ret)
	}
	return
}

func (a *Mutator) DomToMutation(main block.Shape, str string) (err error) {
	//println("dom to mutation", str)
	if els, e := dom.UnmarshalMutation(str); e != nil {
		err = e
	} else if atoms, e := a.mins.LoadMutation(main, a.atomizer, els); e != nil {
		err = e
	} else {
		a.blockPool.CreateMutatedBlock(main, atoms)
	}
	return
}

// create the mui shapes
func (a *Mutator) Decompose(main block.Shape, popup block.Workspace) (block.Shape, error) {
	mb := a.blockPool.EnsureMutatedBlock(main)
	return a.mins.CreateMui(popup, mb)
}

// fill the workspace shape from the mui layout
func (a *Mutator) Compose(main, mui block.Shape) (err error) {
	if target, ok := a.blockPool.GetMutatedBlock(main); !ok {
		err = errutil.New("couldnt find block", main)
	} else if atoms, e := a.mins.DistillMui(target, mui, a.atomizer); e != nil {
		err = e
	} else {
		target.atoms = atoms
	}
	return
}

func (a *Mutator) SaveConnections(main, mui block.Shape) (err error) {
	// note: can be missing the first call.
	if target, ok := a.blockPool.GetMutatedBlock(main); ok {
		target.SaveConnections(mui)
	}
	return
}

func (a *Mutator) PostMixin(b block.Shape) (err error) {
	for _, name := range a.mins.Inputs {
		if in, index := b.InputByName(name); index < 0 {
			err = errutil.New("couldnt find mutation during mixin", b)
		} else {
			in.SetInvisible()
		}
	}
	return
}

func (a *Mutator) Quarks() (ret []string) {
	return PaletteQuarks(a.mins)
}
