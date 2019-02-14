package gblocks

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/named"
)

// collapse all dynamic inputs
func (b *Block) removeAtoms() {
	for i, collapse := 0, 0; i < b.NumInputs(); {
		in := b.Input(i)
		if collapse > 0 {
			b.RemoveInput(in.Name)
			collapse--
		} else {
			if m := in.Mutation(); m != nil {
				collapse = m.ResetAtoms() // returns the total number of inputs used by this mutation
			}
			i++
		}
	}
}

// Build a mutation ui from existing workspace blocks.
// Called once, each time the mutation ui is opened by the user.
// ( as opposed to compose which gets called repeatedly )
func (b *Block) decompose(reg *Registry, mui *Workspace) (ret *Block, err error) {
	if mui == nil {
		err = errutil.New("decompose into nil workspace")
	} else if muiContainer, e := mui.NewBlock(b.MutationType()); e != nil {
		err = e // couldnt get the predefined mutator block used for the dialog
	} else {
		muiContainer.InitSvg()
		// each input in the mui represents an atom in the workspace block
		for i, cnt := 0, muiContainer.NumInputs(); i < cnt; i++ {
			muiInput := muiContainer.Input(i)
			// the input names in the mui match the mutation in the original block
			if blockInput, inputIndex := b.InputByName(muiInput.Name); inputIndex < 0 {
				err = errutil.Append(err, errutil.New("no input named", muiInput.Name))
			} else if m := blockInput.Mutation(); m == nil {
				// ^ the mutation data for the input in the workspace
				err = errutil.Append(err, errutil.New("input isnt mutable", muiInput.Name))
			} else if mutationTypes, ok := reg.mutations.GetMutation(m.MutationName); !ok {
				// ^ the data needed to create mutation ui blocks from atoms
				err = errutil.Append(err, errutil.New("input", muiInput.Name, "has unknown mutation", m.MutationName))
			} else {
				muiConnection := muiInput.Connection()
				// for every atom in the workspace block
				for i, cnt := 0, m.NumAtoms(); i < cnt; i++ {
					atom := m.Atom(i)
					if mutationType, ok := mutationTypes.findMutationType(atom.Type); !ok {
						err = errutil.Append(err, errutil.New("couldnt find mui type for atom", atom.Type))
					} else if muiBlock, e := mui.NewBlock(mutationType); e != nil {
						err = errutil.Append(err, e)
					} else {
						muiBlock.InitSvg()
						muiConnection.Connect(muiBlock.PreviousConnection())
						muiConnection = muiBlock.NextConnection()
					}
				}
			}
		}
		if err == nil {
			ret = muiContainer
		} else {
			muiContainer.Dispose()
		}
	}
	return
}

// before we re/compose the workspace blocks/ remember what the connections pointed to
// connections are saved into the muiContainer's blocks
// when the mui blocks are re-ordered, the connections are re-ordered.
func (b *Block) saveConnections(muiContainer *Block) (err error) {
	// for each input in the mutation ui
	for mi, mcount := 0, muiContainer.NumInputs(); mi < mcount; mi++ {
		muiInput := muiContainer.Input(mi)
		// get the corresponding input in the workspace
		if blockInput, inputIndex := b.InputByName(muiInput.Name); inputIndex < 0 {
			err = errutil.Append(err, errutil.New("no input named", muiInput.Name))
		} else if m := blockInput.Mutation(); m == nil {
			// ^ the mutation data for the input in the workspace
			err = errutil.Append(err, errutil.New("input isnt mutable", muiInput.Name))
		} else {
			atomIndex, atomCnt := 0, m.NumAtoms()
			// visit every block connected to the mutation ui's input
			muiInput.visitStack(func(muiBlock *Block) bool {
				cs := NewConnections()
				// each mutation ui block represents a single atom in the workspace
				// the atom, however, can hold several inputs
				if atomIndex < atomCnt {
					atom := m.Atom(atomIndex)
					atomIndex++
					// record all of the atom's input...
					for i := 0; i < atom.NumInputs; i++ {
						inputIndex++
						in := b.Input(inputIndex)
						cs.AppendInput(in)
					}
				}
				// and store them in the mutation ui's block
				muiBlock.CacheConnections(cs)
				return true // keep going
			})
		}
	}
	return
}

// mutation ui -> workspace
func (b *Block) compose(reg *Registry, muiContainer *Block) (err error) {
	// remove all the dynamic inputs from the blocks;
	// we're about to recreate/recompose them.
	b.removeAtoms()

	var savedInputs []savedMutation
	// for each mutation in the mutator ui
	for mi, mcount := 0, muiContainer.NumInputs(); mi < mcount; mi++ {
		muiInput := muiContainer.Input(mi)
		// get the corresponding input in the workspace
		if blockInput, inputIndex := b.InputByName(muiInput.Name); inputIndex < 0 {
			e := errutil.New("no input named", muiInput.Name)
			err = errutil.Append(err, e)
		} else if m := blockInput.Mutation(); m == nil {
			// ^ the mutation data for the input in the workspace
			e := errutil.New("input isnt mutable", muiInput.Name)
			err = errutil.Append(err, e)
		} else if mutationTypes, ok := reg.mutations.GetMutation(m.MutationName); !ok {
			// ^ the data needed to create workspace atoms from mutation ui blocks
			e := errutil.New("input", muiInput.Name, "has unknown mutation", m.MutationName)
			err = errutil.Append(err, e)
		} else {
			// for every block connected to the mutation ui's input...
			var savedAtoms []savedConnections
			if muiInput.visitStack(func(muiBlock *Block) (keepGoing bool) {
				// determine which workspace atom corresponds to the (user selected) mutation ui block
				if atomType, found := mutationTypes.findAtomType(muiBlock.Type); !found {
					e := errutil.New("unknown atom type for mutation", muiInput.Name, muiBlock.Type)
					err = errutil.Append(err, e)
				} else if numInputs, e := m.addAtom(reg, atomType); e != nil {
					err = errutil.Append(err, e)
				} else if cs := muiBlock.CachedConnections(); cs == nil {
					// can be nil the first time the block gets used
					keepGoing = true
				} else if cnt := cs.Length(); numInputs != cnt {
					e := errutil.New("number of inputs generated by the atom", numInputs, "doesnt match the number of inputs saved for the atom", cs.Length(), muiInput.Name, atomType)
					err = errutil.Append(err, e)
				} else {
					savedAtoms = append(savedAtoms, savedConnections{cs, numInputs})
					keepGoing = true
				}
				return
			}) {
				savedInputs = append(savedInputs, savedMutation{muiInput.Name, savedAtoms})
			}
		}
	}
	//b.InitSvg() -- called by the caller in blockly already.
	if err == nil {
		if e := b.reconnect(savedInputs); e != nil {
			err = e
		} else {
			b.redecorate(reg.Decor)
		}
	}
	return
}

type savedConnections struct {
	connections *Connections
	numInputs   int
}

type savedMutation struct {
	inputName  named.Input
	savedAtoms []savedConnections
}

// re-connect those inputs
func (b *Block) reconnect(savedInputs []savedMutation) (err error) {
	for _, savedInput := range savedInputs {
		inputName := savedInput.inputName
		if in, inputIndex := b.InputByName(inputName); inputIndex < 0 {
			err = errutil.Append(err, errutil.New("no input named", inputName))
		} else if m := in.Mutation(); m == nil {
			err = errutil.Append(err, errutil.New("input isnt mutable", inputName))
		} else {
			for _, savedAtom := range savedInput.savedAtoms {
				cs := savedAtom.connections
				for i, cnt := 0, cs.Length(); i < cnt; i++ {
					c := cs.Connection(i)
					reconnect(b, inputIndex+i+1, c)
				}
				inputIndex += savedAtom.numInputs
			}
		}
	}
	return
}
