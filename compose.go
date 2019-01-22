package gblocks

import (
	"github.com/ionous/errutil"
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

//  workspace -> mutation ui
func (b *Block) decompose(reg *Registry, mui *Workspace) (ret *Block, err error) {
	if mui == nil {
		err = errutil.New("decompose into nil workspace")
	} else if muiContainer, e := mui.NewBlock(b.Type + "$mutation"); e != nil {
		err = e // couldnt get the predefined mutator block used for the dialog
	} else {
		muiContainer.InitSvg()
		// each input in the mui represents an atom in the workspace block
		for i, cnt := 0, muiContainer.NumInputs(); i < cnt; i++ {
			muiInput := muiContainer.Input(i)
			// the input names in the mui match the mutation in the original block
			if blockInput, inputIndex := b.InputByName(muiInput.Name); inputIndex < 0 {
				err = errutil.Append(err, errutil.New("no input named", muiInput.Name))
			} else if m := blockInput.Mutation(); m != nil {
				// ^ the mutation data for the input in the workspace
				err = errutil.Append(err, errutil.New("input isnt mutable", muiInput.Name))
			} else if mutationTypes, ok := reg.mutations[m.MutationName]; !ok {
				// ^ the data needed to create mutation ui blocks from atoms
				err = errutil.Append(err, errutil.New("input", muiInput.Name, "has unknown mutation", m.MutationName))
			} else {
				muiConnection := muiInput.Connection()
				// for every atom in the workspace block
				for i, cnt := 0, m.NumAtoms(); i < cnt; i++ {
					atom := m.Atom(i)
					if mutationType, ok := mutationTypes.findMutationType(atom.Type); !ok {
						err = errutil.Append(err, errutil.New("couldnt type for atom", atom.Type))
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
func (b *Block) saveConnections(muiContainer *Block) (err error) {
	// for each input in the mutation ui
	for mi, mcount := 0, muiContainer.NumInputs(); mi < mcount; mi++ {
		muiInput := muiContainer.Input(mi)
		// get the corresponding input in the workspace
		if blockInput, inputIndex := b.InputByName(muiInput.Name); inputIndex < 0 {
			err = errutil.Append(err, errutil.New("no input named", muiInput.Name))
		} else if m := blockInput.Mutation(); m != nil {
			// ^ the mutation data for the input in the workspace
			err = errutil.Append(err, errutil.New("input isnt mutable", muiInput.Name))
		} else if muiConnection := muiInput.Connection(); muiConnection == nil {
			// ^ the connected blocks in the mutation ui
			err = errutil.Append(err, errutil.New("input is missing connections", muiInput.Name))
		} else {
			// for every block connected to the mutation ui's input...
			atomIndex, atomCnt := 0, m.NumAtoms()
			for muiBlock := muiConnection.TargetBlock(); muiBlock != nil && atomIndex < atomCnt; atomIndex++ {
				// each mutation ui block represents a single atom in the workspace
				atom := m.Atom(atomIndex)
				// the atom, however, can hold several inputs
				connections := NewConnections()
				// record all of the atom's input...
				for i, inputIndex := 0, inputIndex+1; i < atom.NumInputs; i, inputIndex = i+1, inputIndex+1 {
					in := b.Input(inputIndex)
					connections.Append(in.Connection().TargetConnection())
				}
				// and store them in the mutation ui's block
				muiBlock.connections = connections
				// then, move to the next block connected to the mutation ui's input
				if muiConnection := muiBlock.NextConnection(); muiConnection != nil {
					muiBlock = muiConnection.TargetBlock()
				} else {
					break // done with blocks
				}
			}
		}
	}
	return
}

// mutation ui -> workspace
func (b *Block) compose(reg *Registry, muiContainer *Block) (err error) {
	// remove all the dynamic inputs from the blocks;
	// we're about to recreate/recompose them.
	b.removeAtoms()
	type savedConnection struct {
		connections *Connections
		numInputs   int
	}
	type savedMutation struct {
		inputName        InputName
		savedConnections []savedConnection
	}
	var savedMutations []savedMutation
	// for each mutation in the mutator ui
	for mi, mcount := 0, muiContainer.NumInputs(); mi < mcount; mi++ {
		muiInput := muiContainer.Input(mi)
		// get the corresponding input in the workspace
		if blockInput, inputIndex := b.InputByName(muiInput.Name); inputIndex < 0 {
			err = errutil.Append(err, errutil.New("no input named", muiInput.Name))
		} else if m := blockInput.Mutation(); m != nil {
			// ^ the mutation data for the input in the workspace
			err = errutil.Append(err, errutil.New("input isnt mutable", muiInput.Name))
		} else if mutationTypes, ok := reg.mutations[m.MutationName]; !ok {
			// ^ the data needed to create workspace atoms from mutation ui blocks
			err = errutil.Append(err, errutil.New("input", muiInput.Name, "has unknown mutation", m.MutationName))
		} else {
			var savedConnections []savedConnection
			if muiConnection := muiInput.Connection(); muiConnection != nil {
				okay := true
				// for every block connected to the mutation ui's input...
				for muiBlock := muiConnection.TargetBlock(); muiBlock != nil; {
					// determine which workspace atom corresponds to the (user selected) mutation ui block
					if atomType, found := mutationTypes.findAtomType(muiBlock.Type); !found {
						err = errutil.Append(err, errutil.New("unknown atom type for mutation", muiInput.Name, muiBlock.Type))
						okay = false
						break
					} else if numInputs, e := m.addAtom(reg, atomType); e != nil {
						err = errutil.Append(err, e)
						okay = false
						break
					} else {
						// after we have generated all the inputs in the workspace block we will need to reconnect them.
						saved := savedConnection{muiBlock.connections, numInputs}
						savedConnections = append(savedConnections, saved)

						// next clause in the mutation ui ( for this input )
						if muiConnection := muiBlock.NextConnection(); muiConnection != nil {
							muiBlock = muiConnection.TargetBlock()
						} else {
							break
						}
					}
					if okay {
						savedMutations = append(savedMutations, savedMutation{muiInput.Name, savedConnections})
					}
				}
			}
			b.InitSvg()
			// re-connect those inputs
			for _, saved := range savedMutations {
				name := saved.inputName
				if in, index := b.InputByName(name); index < 0 {
					err = errutil.Append(err, errutil.New("missing input named", name))
				} else if m := in.Mutation(); m != nil {
					err = errutil.Append(err, errutil.New("missing mutable data", name))
				} else {
					for _, saved := range saved.savedConnections {
						inputs, cs := saved.numInputs, saved.connections
						var cnt int
						if a, b := inputs, cs.Length(); a < b {
							cnt = a
						} else {
							cnt = b
						}
						for i := 0; i < cnt; i++ {
							reconnect(b, index+i+1, cs.Connection(i))
						}
					}
				}
			}
		}
	}
	return
}
