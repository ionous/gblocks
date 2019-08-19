// Package blockly wraps google's blocky javascript api for use with gopherjs.
// see also:
// https://developers.google.com/blockly
// https://github.com/gopherjs/gopherjs
package blockly

import "github.com/gopherjs/gopherjs/js"

type Blockly struct {
	*js.Object
	blocks     *js.Object `js:"Blocks"`
	xml        *js.Object `js:"Xml"`
	extensions *js.Object `js:"Extensions"`
	utils      *js.Object `js:"utils"`
}
