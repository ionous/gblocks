package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	r "reflect"
)

type Mutation interface {
	// return the value of the internal array
	// this relieves the user from implementing element access ( num, set, get )
	Elements() r.Value
	MutationForType(r.Type) r.Type
}

func reconnect(connectionChild *Connection, block *Block, inputName string) (okay bool) {
	if res := js.Global.Get("Blockly").Get("Mutator").Call("reconnect",
		connectionChild.Object, block.Object, inputName); res.Bool() {
		okay = true
	}
	return
}

// get the users's first placed block from the mutation dialog.
// either its the one after the "head block", or its the one inside a c-shape statement input
func firstConnection(containerBlock *Block) (ret *Connection) {
	if next := containerBlock.NextConnection; next != nil {
		ret = next
	} else if inner := containerBlock.GetFirstStatementConnection(); inner != nil {
		ret = inner
	}
	return
}

func firstClause(containerBlock *Block) (ret *Block) {
	if c := firstConnection(containerBlock); c != nil {
		ret = c.TargetBlock()
	}
	return
}
