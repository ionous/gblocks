package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/named"
	r "reflect"
)

// Mutator -- blockly api
type Mutator struct {
	*js.Object
}

func NewMutator(quarkNames []named.Type) (ret *Mutator) {
	if blockly := GetBlockly(); blockly != nil {
		obj := blockly.Get("Mutator").New(quarkNames)
		ret = &Mutator{Object: obj}
	}
	return
}

// Mutation - user specification of a mutation block.
type Mutation struct {
	Label   string
	Creates interface{}
}

// MutationBlock - gblocks internal description for the palette of a mutation ui popup
// the mutation ui blocks are auto-generated and limited in appearance:
// a name, a generic previous connection, a possible next connection.
type MutationBlock struct {
	MuiLabel      string // used as the label for the block in the ui
	MuiType       named.Type
	WorkspaceType named.Type // type of the top block created by block xml; same as Xml["type"]
	Constraints              // mutation ui block types permitted to follow this block.
	//BlockXml      *XmlElement // workspace block xml duplicated the when the mutaiton block gets newly placed.
}

// RegisteredMutation - gblocks internal description of the palette used by a mutation popup.
type RegisteredMutation struct {
	// MuiType -> MutationBlock
	blocks map[named.Type]*MutationBlock
	quarks []named.Type // keys of blocks in display order.
}

// RegisteredMutations -
type RegisteredMutations struct {
	typeToMutation map[named.Type]*RegisteredMutation
}

/**
 * Reconnect an block to a mutated input.
 * @return {boolean} True iff a reconnection was made, false otherwise.
 */
func reconnect(block *Block, i int, tgtConnection *Connection) (okay bool) {
	if tgtConnection != nil {
		// ensure the block hasnt been disposed.
		if parentBlock := tgtConnection.GetSourceBlock(); parentBlock != nil && parentBlock.hasWorkspace() {
			if src := block.Input(i).Connection(); src != nil {
				targetBlock := tgtConnection.TargetBlock()
				if ((targetBlock == nil) || (targetBlock == block)) && src.TargetConnection() != tgtConnection {
					if src.IsConnected() {
						// There's already something connected here.  Get rid of it.
						src.Disconnect()
					}
					src.Connect(tgtConnection)
					okay = true
				}
			}
		}
	}
	return
}

//
func (m *MutationBlock) BlockDesc() Dict {
	blockDesc := Dict{
		opt_type:     m.MuiType,
		opt_message0: m.MuiLabel,
		opt_previous: nil, // all have a generic previous connection
	}
	if types, constrained := m.GetConstraints(); constrained {
		blockDesc[opt_next] = types
	}
	return blockDesc
}

func (m *MutationBlock) BlockFns() Dict {
	blockDesc := m.BlockDesc()
	init := js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
		b := &Block{Object: obj}
		if e := b.JsonInit(blockDesc); e != nil {
			panic(e)
		}
		return
	})
	return Dict{
		"init": init,
	}
}

func (reg *RegisteredMutations) Contains(mutation named.Type) (okay bool) {
	if _, ok := reg.typeToMutation[mutation]; ok {
		okay = true
	}
	return
}

func (reg *RegisteredMutations) GetMutation(mutation named.Type) (ret *RegisteredMutation, okay bool) {
	if r, ok := reg.typeToMutation[mutation]; ok {
		ret, okay = r, true
	}
	return
}

func (reg *RegisteredMutations) RegisterMutation(mutation named.Type, muiBlocks ...Mutation) (err error) {
	if reg.Contains(mutation) {
		err = errutil.New("mutation already exists", mutation)
	} else if blockly := GetBlockly(); blockly == nil {
		err = errutil.New("blockly doesnt exist")
	} else {
		// add the "tops" of the prototypes to the pool we pull from to connect "next" blocks.
		types := make(RegisteredTypes)
		for _, el := range muiBlocks {
			structType := r.TypeOf(el.Creates).Elem()
			types.RegisterType(structType)
		}

		// now, walk those prototypes again
		var quarks []named.Type
		blocks := make(map[named.Type]*MutationBlock)

		for _, muiBlock := range muiBlocks {
			prototype, label := muiBlock.Creates, muiBlock.Label
			structType := r.TypeOf(prototype).Elem()
			typeName := named.TypeFromStruct(structType)

			// scan to the end of the prototype's NextStatement stack
			// lastVal := val
			// var lastField []int
			// for {
			// 	lastType := lastVal.Type()
			// 	if f, ok := lastType.FieldByName(NextField); !ok {
			// 		lastField = nil // there is no next field; clear anything from a previous block in the chain
			// 		break
			// 	} else if nextVal := val.FieldByIndex(f.Index); !nextVal.IsValid() || nextVal.IsNil() {
			// 		lastField = f.Index
			// 		break
			// 	} else {
			// 		lastVal = nextVal.Elem()
			// 		lastType = nextVal.Type()
			// 	}
			// }
			// var constraints Constraints
			// if len(lastField) > 0 {
			// 	if c, e := types.CheckStructField(lastVal.Type().FieldByIndex(lastField)); e != nil {
			// 		err = errutil.Append(err, e)
			// 	} else {
			// 		constraints = c
			// 	}
			// }

			// future: prototype into dom tree
			// xml := ValueToDom(v, true)
			// // does the element have sub-elements (or is it just one block?)
			// var subElements int
			// if shadows := xml.GetElementsByTagName("shadow"); shadows != nil {
			// 	subElements = shadows.Num()
			// }

			if constraints, e := types.CheckField(structType, NextField); e != nil {
				err = errutil.Append(err, e)
			} else {
				muiType := named.SpecialType("mui", mutation.String(), typeName.String())
				muiBlock := &MutationBlock{
					MuiLabel:      label,
					MuiType:       muiType,
					WorkspaceType: typeName,
					Constraints:   constraints,
					//		BlockXml:      xml,
				}
				blockly.AddBlock(muiType, muiBlock.BlockFns())
				quarks = append(quarks, muiType)
				blocks[muiType] = muiBlock
			}

			//
			// if !dupes[typeName] {
			// 	// if there are sub-elements; then we have also register the first block
			// 	if isAtom := subElements == 0; !isAtom {
			// 		muiType := named.SpecialType"mui", name, el.Name, "atom")
			// 		b := &MutationBlock{
			// 			MuiLabel: typeName.Friendly(),
			// 			WorkspaceType: typeName,
			// 			BlockXml:      xml,
			// 			// FIXXXX -- these constraints are probably wrong...
			// 			Constraints: constraints,
			// 		}
			// 		blockly.AddBlock(muiType, b.BlockFns())
			// 		atoms = append(atoms, muiType)

			// 	}
			// 	// regardless, we either have added the atom, or the block itself was an atom.
			// 	dupes[typeName] = true
			// }

		}
		// append the atoms at the end of the other blocks
		// quarkNames, atomNames = append(quarkNames, atomNames...), nil

		// TODO: color code the blocks by the mui's container input -- each input a different set of shades.
		if reg.typeToMutation == nil {
			reg.typeToMutation = make(map[named.Type]*RegisteredMutation)
		}
		reg.typeToMutation[mutation] = &RegisteredMutation{blocks: blocks, quarks: quarks}
	}
	return
}

// what kind of workspace block/s the mutation element represent*
func (mm *RegisteredMutation) findMutationType(wsType named.Type) (ret named.Type, okay bool) {
	for muiType, m := range mm.blocks {
		if m.WorkspaceType == wsType {
			ret = muiType
			okay = true
			break
		}
	}
	return
}

//  what kind of workspace block/s the mutation element represent*
func (mm *RegisteredMutation) findAtomType(muiType named.Type) (ret named.Type, okay bool) {
	if m, ok := mm.blocks[muiType]; ok {
		ret = m.WorkspaceType
		okay = true
	}
	return
}
