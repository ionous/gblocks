package mutant

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
)

// instance data for mutations
// implements block.Mutator
type Mutator struct {
	mins          *BlockMutations // descriptions of atoms
	atomizer      Atomizer        // needed to expand quarks into atoms
	mutableBlocks MutatedBlocks   // per-block mutation data (pooled across mutators)
}

func NewMutator(mins *BlockMutations, db Atomizer, mbs MutatedBlocks) *Mutator {
	return &Mutator{mins, db, mbs}
}

// pool of user directed block mutations.
type MutatedBlocks map[string]*mutableBlock

func (mbs MutatedBlocks) OnDelete(id string) {
	delete(mbs, id)
}

type mutableBlock struct {
	inputs      MutableInputs
	connections SavedConnections
}

func (a *Mutator) Delete(id string) {
	delete(a.mutableBlocks, id)
}

func (a *Mutator) MutationToDom(main block.Shape) (ret string, err error) {
	id := main.BlockId()
	if m, ok := a.mutableBlocks[id]; ok {
		dom := a.mins.SaveMutation(m.inputs)
		ret, err = dom.MarshalMutation()
		//println("mutation to dom", ret)
	}
	return
}

func (a *Mutator) DomToMutation(main block.Shape, str string) (err error) {
	//println("dom to mutation", str)
	if els, e := dom.UnmarshalMutation(str); e != nil {
		err = e
	} else if inputs, e := a.mins.LoadMutation(main, a.atomizer, els); e != nil {
		err = e
	} else {
		id := main.BlockId()
		a.mutableBlocks[id] = &mutableBlock{inputs: inputs}
	}
	return
}

// create the mui shapes
func (a *Mutator) Decompose(main block.Shape, popup block.Workspace) (block.Shape, error) {
	id := main.BlockId()
	// does this block have previous mutation data?
	var inputs MutableInputs
	if m, ok := a.mutableBlocks[id]; ok {
		inputs = m.inputs
	} else {
		// no? it's a decent time to create it.
		a.mutableBlocks[id] = new(mutableBlock)
	}
	return a.mins.CreateMui(popup, main, inputs)
}

// fill the workspace shape from the mui layout
func (a *Mutator) Compose(main, mui block.Shape) (err error) {
	id := main.BlockId()
	if m, ok := a.mutableBlocks[id]; !ok {
		err = errutil.New("couldnt find block", id)
	} else if inputs, e := a.mins.DistillMui(main, mui, a.atomizer, m.connections); e != nil {
		err = e
	} else {
		m.inputs = inputs
	}
	return
}

func (a *Mutator) SaveConnections(main, mui block.Shape) (err error) {
	// note: can be missing the first call.
	id := main.BlockId()
	if m, ok := a.mutableBlocks[id]; ok {
		cs := SaveConnections(main, mui)
		m.connections = cs
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
