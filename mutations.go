package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type Mutator struct {
	*js.Object
}

func NewMutator(quarkNames []TypeName) (ret *Mutator) {
	if blockly := js.Global.Get("Blockly"); blockly.Bool() {
		obj := blockly.Get("Mutator").New(quarkNames)
		ret = &Mutator{Object: obj}
	}
	return
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

// Mutation - user specification of a mutation block.
type Mutation struct {
	Label   string
	Creates interface{}
}

// MutationBlock - internal description for the palette of a mutation ui popup
// the mutation ui blocks are auto-generated and limited in appearance:
// a name, a generic previous connection, a possible next connection.
type MutationBlock struct {
	MuiLabel      string // used as the label for the block in the ui
	MuiType       TypeName
	WorkspaceType TypeName // type of the top block created by block xml; same as Xml["type"]
	Constraints            // mutation ui block types permitted to follow this block.
	//BlockXml      *XmlElement // workspace block xml duplicated the when the mutaiton block gets newly placed.
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

// RegisteredMutation - internal description of the palette used by a mutation popup.
type RegisteredMutation struct {
	// MuiType -> MutationBlock
	blocks map[TypeName]*MutationBlock
	quarks []TypeName // keys of blocks in display order.
}

type RegisteredMutations map[string]*RegisteredMutation

// what kind of workspace block/s the mutation element represent*
func (mm *RegisteredMutation) findMutationType(wsType TypeName) (ret TypeName, okay bool) {
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
func (mm *RegisteredMutation) findAtomType(muiType TypeName) (ret TypeName, okay bool) {
	if m, ok := mm.blocks[muiType]; ok {
		ret = m.WorkspaceType
		okay = true
	}
	return
}

// could be a map, except maps arent ordered.
type mutationInput struct {
	inputName    InputName //
	mutationName string    // the name as per RegisterMutation, and the struct tags
	constraints  Constraints
}
