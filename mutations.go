package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

func reconnect(connectionChild *Connection, block *Block, inputName string) (okay bool) {
	if res := js.Global.Get("Blockly").Get("Mutator").Call("reconnect",
		connectionChild.Object, block.Object, inputName); res.Bool() {
		okay = true
	}
	return
}
