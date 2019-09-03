package mutant

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
	"github.com/kr/pretty"
)

// instance data for mutations
// implements block.Mutator
type Mutator struct {
	arch     *BlockMutations // descriptions of atoms
	atomizer Atomizer        // needed to expand quarks into atoms
	// everyone has the same atomizer, so consider putting in the mutated blocks
	blockPool *MutatedBlocks // per-block mutation data (pooled across mutators)
}

func NewMutator(arch *BlockMutations, db Atomizer, mbs *MutatedBlocks) *Mutator {
	return &Mutator{arch, db, mbs}
}

var Trace bool = true

func (a *Mutator) MutationToDom(main block.Shape) (ret string, err error) {
	if Trace {
		println("saving mutation")
	}
	if src, ok := a.blockPool.GetMutatedBlock(main); ok {
		dom := src.SaveMutation()
		ret, err = dom.MarshalMutation()
	}
	return
}

func (a *Mutator) DomToMutation(main block.Shape, str string) (err error) {
	if els, e := dom.UnmarshalMutation(str); e != nil {
		err = e
	} else {
		if Trace {
			pretty.Println("loading mutation", els)

		}
		target := a.blockPool.CreateMutatedBlock(main, a.arch, a.atomizer)
		err = target.LoadMutation(&els)
	}
	return
}

// create the mui container from workspace data.
func (a *Mutator) Decompose(main block.Shape, popup block.Workspace) (block.Shape, error) {
	if Trace {
		println("create mui")
	}
	src := a.blockPool.CreateMutatedBlock(main, a.arch, a.atomizer)
	return src.CreateMui(popup)
}

// create workspace inputs from the atoms the user selected and arranged in the mui popup
func (a *Mutator) Compose(main, mui block.Shape) (err error) {
	if Trace {
		println("create from mui")
	}
	if target, ok := a.blockPool.GetMutatedBlock(main); !ok {
		err = errutil.New("couldnt find block", main)
	} else if e := target.CreateFromMui(mui); e != nil {
		err = e
	}
	return
}

func (a *Mutator) SaveConnections(main, mui block.Shape) (err error) {
	// note: can be missing the first call.
	if target, ok := a.blockPool.GetMutatedBlock(main); ok {
		res := target.SaveConnections(mui)
		if Trace {
			println("saved connections:")
			for x, store := range res {
				println(" atom", x)
				for _, c := range store.Connections() {
					type stringer interface{ String() string }
					if stringer, ok := c.(stringer); ok {
						println("  ", stringer.String())
					}
				}
			}
		}
	}
	return
}

func (a *Mutator) PostMixin(b block.Shape) (err error) {
	for _, name := range a.arch.Inputs {
		if in, index := b.InputByName(name); index < 0 {
			err = errutil.New("couldnt find mutation during mixin", b)
		} else {
			in.SetInvisible()
		}
	}
	return
}

func (a *Mutator) Quarks() (ret []string) {
	return PaletteQuarks(a.arch)
}
