package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/named"
)

type Blockly struct {
	*js.Object
	blocks *js.Object `js:"Blocks"`
	xml    *js.Object `js:"Xml"`
}

func GetBlockly() (ret *Blockly) {
	if obj := js.Global.Get("Blockly"); obj.Bool() {
		ret = &Blockly{Object: obj}
	}
	return ret
}

func (b *Blockly) AddBlock(typeName named.Type, fns Dict) {
	b.blocks.Set(typeName.String(), fns)
}

func (b *Blockly) Xml() *Xml {
	return &Xml{Object: b.xml}
}
