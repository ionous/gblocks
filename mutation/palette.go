package mutation

import (
	// 	"github.com/gopherjs/gopherjs/js"
	// 	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	r "reflect"
)

// // Mutation - user specification of a mutation block.
// type Mutation struct {
// 	Label   string
// 	Creates interface{}
// }

// Palette - description of the palette used by a mutation popup.
type Palette struct {
	mutationType r.Type                       // name of mutation type
	entries      map[block.Type]*PaletteEntry // MuiType -> PaletteEntry
	quarks       []block.Type                 // keys of palette in display order.
}

// PaletteEntry - palette of a mutation ui popup
// the mutation ui blocks are auto-generated and limited in appearance:
// a name, a generic previous connection, a possible next connection.
type PaletteEntry struct {
	MuiLabel      string     // used as the label for the block in the mui
	MuiType       block.Type // name of the mutation block in the mui
	WorkspaceType block.Type // name of the workspace block that the mutation block represents
	//MuiConstraints            // mutation ui block types permitted to follow this block.
	//BlockXml      *Element // workspace block xml duplicated the when the mutaiton block gets newly placed.
}

// what kind of workspace block/s the mutation element represent*
func (p *Palette) findMutationType(wsType block.Type) (ret block.Type, okay bool) {
	for muiType, m := range p.entries {
		if m.WorkspaceType == wsType {
			ret, okay = muiType, true
			break
		}
	}
	return
}

//  what kind of workspace block/s the mutation element represent*
func (p *Palette) findAtomType(muiType block.Type) (ret block.Type, okay bool) {
	if m, ok := p.entries[muiType]; ok {
		ret, okay = m.WorkspaceType, true
	}
	return
}

// //
// func (m *PaletteEntry) BlockDesc() Dict {
// 	blockDesc := Dict{
// 		OptType:     m.MuiType,
// 		OptMessage0: m.MuiLabel,
// 		OptPrevious: nil, // all have a generic previous connection
// 	}
// 	if types, constrained := m.GetConstraints(); constrained {
// 		blockDesc[OptNext] = types
// 	}
// 	return blockDesc
// }

// func (m *PaletteEntry) BlockFns() Dict {
// 	blockDesc := m.BlockDesc()
// 	init := js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
// 		b := &Block{Object: obj}
// 		if e := b.JsonInit(blockDesc); e != nil {
// 			panic(e)
// 		}
// 		return
// 	})
// 	return Dict{
// 		"init": init,
// 	}
// }

// func NewPalette(mutationType r.Type, muiBlocks ...Mutation) (ret *Palette, err error) {
// 	mutation := block.TypeFromStruct(mutationType)
// 	if reg.Contains(mutation) {
// 		err = errutil.New("mutation already exists", mutation)
// 	} else if blockly := GetBlockly(); blockly == nil {
// 		err = errutil.New("blockly doesnt exist")
// 	} else {
// 		// add the "tops" of the prototypes to the pool we pull from to connect "next" blocks.
// 		types := make(RegisteredTypes)
// 		for _, el := range muiBlocks {
// 			structType := r.TypeOf(el.Creates).Elem()
// 			types.RegisterType(structType)
// 		}

// 		// now, walk those prototypes again
// 		var quarks []block.Type
// 		blocks := make(map[block.Type]*PaletteEntry)

// 		for _, muiBlock := range muiBlocks {
// 			prototype, label := muiBlock.Creates, muiBlock.Label
// 			structType := r.TypeOf(prototype).Elem()
// 			typeName := block.TypeFromStruct(structType)

// 			// scan to the end of the prototype's NextStatement stack
// 			// lastVal := val
// 			// var lastField []int
// 			// for {
// 			// 	lastType := lastVal.Type()
// 			// 	if f, ok := lastType.FieldByName(NextStatement); !ok {
// 			// 		lastField = nil // there is no next field; clear anything from a previous block in the chain
// 			// 		break
// 			// 	} else if nextVal := val.FieldByIndex(f.Index); !nextVal.IsValid() || nextVal.IsNil() {
// 			// 		lastField = f.Index
// 			// 		break
// 			// 	} else {
// 			// 		lastVal = nextVal.Elem()
// 			// 		lastType = nextVal.Type()
// 			// 	}
// 			// }
// 			// var constraints Constraints
// 			// if len(lastField) > 0 {
// 			// 	if c, e := types.CheckStructField(lastVal.Type().FieldByIndex(lastField)); e != nil {
// 			// 		err = errutil.Append(err, e)
// 			// 	} else {
// 			// 		constraints = c
// 			// 	}
// 			// }

// 			// future: prototype into dom tree
// 			// xml := ValueToDom(v, true)
// 			// // does the element have sub-elements (or is it just one block?)
// 			// var subElements int
// 			// if shadows := xml.GetElementsByTagName("shadow"); shadows != nil {
// 			// 	subElements = shadows.Num()
// 			// }

// 			if constraints, e := types.CheckField(structType, StatementNext); e != nil {
// 				err = errutil.Append(err, e)
// 			} else {
// 				muiType := block.SpecialType("mui", mutation.String(), typeName.String())
// 				muiBlock := &PaletteEntry{
// 					MuiLabel:      label,
// 					MuiType:       muiType,
// 					WorkspaceType: typeName,
// 					Constraints:   constraints,
// 					//		BlockXml:      xml,
// 				}
// 				blockly.AddBlock(muiType, muiBlock.BlockFns())
// 				quarks = append(quarks, muiType)
// 				blocks[muiType] = muiBlock
// 			}

// 			//
// 			// if !dupes[typeName] {
// 			// 	// if there are sub-elements; then we have also register the first block
// 			// 	if isAtom := subElements == 0; !isAtom {
// 			// 		muiType := block.SpecialType"mui", name, el.Name, "atom")
// 			// 		b := &PaletteEntry{
// 			// 			MuiLabel: typeName.Friendly(),
// 			// 			WorkspaceType: typeName,
// 			// 			BlockXml:      xml,
// 			// 			// FIXXXX -- these constraints are probably wrong...
// 			// 			Constraints: constraints,
// 			// 		}
// 			// 		blockly.AddBlock(muiType, b.BlockFns())
// 			// 		atoms = append(atoms, muiType)

// 			// 	}
// 			// 	// regardless, we either have added the atom, or the block itself was an atom.
// 			// 	dupes[typeName] = true
// 			// }

// 		}
// 		// append the atoms at the end of the other blocks
// 		// quarkNames, atomNames = append(quarkNames, atomNames...), nil

// 		// TODO: color code the blocks by the mui's container input -- each input a different set of shades.

// 		newMutation := &Palette{mutationType: mutationType, blocks: blocks, quarks: quarks}
// 	}
// 	return
// }
