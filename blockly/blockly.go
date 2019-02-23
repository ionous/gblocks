// Package blockly wraps google's blocky javascript api for use with gopherjs.
// see also:
// https://developers.google.com/blockly
// https://github.com/gopherjs/gopherjs
package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/block"
)

type Blockly struct {
	*js.Object
	blocks     *js.Object `js:"Blocks"`
	xml        *js.Object `js:"Xml"`
	extensions *js.Object `js:"Extensions"`
}

// blockly is a global object
func getBlockly() (ret *Blockly) {
	if obj := js.Global.Get("Blockly"); obj.Bool() {
		ret = &Blockly{Object: obj}
	}
	return ret
}

// func AddBlock(typeName block.Type, fns block.Dict) {
// 	b := getBlockly()
// 	b.blocks.Set(typeName.String(), fns)
// }

func DefineBlock(typeName block.Type, blockDesc block.Dict) {
	blockly := getBlockly()
	blockly.Call("jsonInitFactory_", blockDesc)
}

func Get(name string) *js.Object {
	return getBlockly().Get(name)
}
