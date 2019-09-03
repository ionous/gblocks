package mutant

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
)

// each mutable workspace block contains
// one or more mutable inputs containing a list of atoms
// and connections to existing workspace blocks.
type MutatedBlock struct {
	block    block.Shape     // workspace block
	arch     *BlockMutations // data about the mutation
	atomizer Atomizer
	inputs   map[string]*MutatedInput
	// muiblock id / atom id to array of connections
	// note: the store cant live with the atom b/c the atom can get destroyed
	store map[string]Store
}

func (mb *MutatedBlock) ContainerName() string {
	return ContainerName(mb.block.BlockType())
}

// returns a MutatedInput if the named input was mutable.
func (mb *MutatedBlock) GetMutatedInput(inputName string) (ret *MutatedInput, okay bool) {
	if a, ok := mb.inputs[inputName]; ok {
		ret, okay = a, true
	} else if min, ok := mb.arch.GetMutation(inputName); ok {
		// verify the input is mutable before creating it.
		a = &MutatedInput{InputName: inputName, Arch: min}
		mb.inputs[inputName] = a
		ret, okay = a, true
	}
	return
}

// SaveConnections happens before any mui transformation
func (mb *MutatedBlock) SaveConnections(muiContainer block.Shape) Storage {
	// 1. walk mutable inputs
	storage := make(Storage)
	for mi, mcnt := 0, muiContainer.NumInputs(); mi < mcnt; mi++ {
		muiInput := muiContainer.Input(mi)
		// 2. walk the mui blocks attached to the mutable input
		block.VisitStack(muiInput, func(muiBlock block.Shape) (keepGoing bool) {
			var store Store
			// 3. get mui block name; this is the unique id of an atom
			atomId := muiBlock.BlockId()
			// 4. find ws inputs with that name
			scope := block.Scope("a", atomId)
			for bi, bcnt := 0, mb.block.NumInputs(); bi < bcnt; bi++ {
				in := mb.block.Input(bi)
				if strings.HasPrefix(in.InputName(), scope) {
					// 5. store outgoing connections
					store.SaveConnection(in)
				}
			}
			storage[atomId] = store
			return true
		})
	}
	mb.store = storage
	return storage
}

// CreateMui from existing workspace blocks. ( aka decompose )
func (mb *MutatedBlock) CreateMui(mui block.Workspace) (ret block.Shape, err error) {
	containerType := mb.ContainerName()
	if muiContainer, e := mui.NewBlock(containerType); e != nil {
		err = e
	} else {
		muiContainer.InitSvg() // from blockly examples; to finalize the inputs?
		for _, inputName := range mb.arch.Inputs {
			if min, ok := mb.GetMutatedInput(inputName); !ok {
				e := errutil.New("input not mutable", inputName)
				err = errutil.Append(err, e)
			} else if muiInput, dex := muiContainer.InputByName(inputName); dex < 0 {
				e := errutil.New("can't find container input", inputName)
				err = errutil.Append(err, e)
			} else {
				stack := muiInput.Connection()
				// walk existing atoms
				for _, atom := range min.Atoms {
					// create blocks named after the atom
					if q, ok := FindQuark(min.Arch, atom.Type); !ok {
						e := errutil.New("couldnt find quark for atom", atom.Type)
						err = errutil.New(err, e)
					} else if muiBlock, e := mui.NewBlockWithId(atom.Name, q.BlockType()); e != nil {
						err = errutil.Append(err, e)
					} else {
						// link the new block into the stack
						muiBlock.InitSvg()
						stack.Connect(muiBlock.PreviousConnection())
						stack = muiBlock.NextConnection()
					}
				}
			}
		}
		if err != nil {
			muiContainer.Dispose()
		} else {
			ret = muiContainer
		}
	}
	return
}

// aka. compose -- turn the mui into new workspace inputs
// adds new inputs to target, returns the atoms for those inputs
func (mb *MutatedBlock) CreateFromMui(muiContainer block.Shape) (err error) {
	// remove all the dynamic inputs from the blocks; we're about to recreate/recompose them.
	// note: the connections for those inputs have already been saved.
	RemoveAtoms(mb.block)
	// create atoms from mui blocks
	for _, inputName := range mb.arch.Inputs {
		if min, ok := mb.GetMutatedInput(inputName); !ok {
			e := errutil.New("input not mutable", inputName)
			err = errutil.Append(err, e)
		} else if muiInput, dex := muiContainer.InputByName(inputName); dex < 0 {
			e := errutil.New("can't find container input", inputName)
			err = errutil.Append(err, e)
		} else {
			var atoms []*AtomizedInput
			// walk the mui blocks attached to the mutable input
			block.VisitStack(muiInput, func(muiBlock block.Shape) (keepGoing bool) {
				// get mui block name; this is the unique id of an atom
				// FIX? consider storing the blockId as atomType, atomId instead of searching by quark
				// FIX? better, why cant the block type *be* the atom type?
				atomId := muiBlock.BlockId()
				blockType := muiBlock.BlockType()
				// translate the block type into atom Type
				if q, ok := FindQuark(min.Arch, blockType); !ok {
					e := errutil.New("can't find quark", blockType)
					err = errutil.Append(err, e)
				} else {
					atomType := q.Name()
					atoms = append(atoms, &AtomizedInput{
						Name: atomId,
						Type: atomType,
					})
				}
				return true
			})
			min.Atoms = atoms
		}
	}
	// make ws inputs from atoms
	if e := mb.expandAtoms(); e != nil {
		err = errutil.Append(err, e)
	}
	return
}

// LoadMutation, creating workspace inputs based on the recorded mutation
func (mb *MutatedBlock) LoadMutation(els *dom.BlockMutation) (err error) {
	// read mutation into atoms
	for _, el := range els.Inputs {
		// watch for empty mutations ( where no input data exists )
		if inputName := el.Input; len(inputName) > 0 {
			if min, ok := mb.GetMutatedInput(inputName); !ok {
				e := errutil.New("input not mutable", inputName)
				err = errutil.Append(err, e)
			} else {
				var atoms []*AtomizedInput
				for _, atom := range el.Atoms {
					atoms = append(atoms, &AtomizedInput{
						Name: atom.Name,
						Type: atom.Type,
					})
				}
				min.Atoms = atoms
			}
		}
	}
	// we're refreshing from xml, no connections need saving
	mb.store = nil
	// make ws inputs from atoms
	if e := mb.expandAtoms(); e != nil {
		err = errutil.Append(err, e)
	}
	return
}

// serialize the mutation to xml friendly data
func (mb *MutatedBlock) SaveMutation() (ret dom.BlockMutation) {
	for inputName, min := range mb.inputs {
		// if there are atoms, create a node for the data.
		var out []*dom.Atom
		for _, atom := range min.Atoms {
			diatom := &dom.Atom{Name: atom.Name, Type: atom.Type}
			out = append(out, diatom)
		}
		if len(out) != 0 {
			ret.Append(&dom.Mutation{inputName, out})
		}
	}
	return
}

// make ws inputs from atoms
func (mb *MutatedBlock) expandAtoms() (err error) {
	// remove all the dynamic inputs from the blocks; we're about to recreate/recompose them.
	// note: the connections for those inputs have already been saved.
	RemoveAtoms(mb.block)
	//
	for _, min := range mb.inputs {
		if j, e := newInjector(mb, min); e != nil {
			err = errutil.Append(err, e)
		} else {
			//
			for _, atom := range min.Atoms {
				if q, ok := FindQuark(min.Arch, atom.Type); !ok {
					e := errutil.New("unknown atom", atom)
					err = errutil.Append(err, e)
				} else {
					scope := block.Scope("a", atom.Name)
					if args, e := q.Atomize(scope, mb.atomizer); e != nil {
						err = errutil.Append(err, e)
					} else {
						// args appear at the end
						if start, cnt := j.inject(args); cnt > 0 {
							if store, ok := mb.store[atom.Name]; ok {
								store.Restore(mb.block, start, cnt)
							}
						}
					}
				}
			}
			// sort the inputs which were added into the right spot.
			j.finalizeInputs()
		}
	}
	return // err
}
