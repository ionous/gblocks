package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type Blockly struct {
	*js.Object
	blocks *js.Object `js:"Blocks"`
}

func GetBlockly() (ret *Blockly) {
	if obj := js.Global.Get("Blockly"); obj.Bool() {
		ret = &Blockly{Object: obj}
	}
	return ret
}

func (b *Blockly) AddBlock(typeName TypeName, fns Dict) {
	b.blocks.Set(typeName.String(), fns)
}
