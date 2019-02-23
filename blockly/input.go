package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

// InputType - describes both inputs and connections
type InputType int

//go:generate stringer -type=InputType
const (
	InputValueType InputType = iota + 1
	OutputValueType
	NextStatementType     // used for connections between blocks, and for statement inputs
	PreviousStatementType // used for connections between blocks
	DummyInputType
)

type InputAlign int

//go:generate stringer -type=InputAlign
const (
	AlignLeft InputAlign = iota - 1
	AlignCentre
	AlignRight
)

type Input struct {
	*js.Object             // Blockly.Input
	Type        InputType  `js:"type"`
	Name        block.Item `js:"name"`
	Align       InputAlign `js:"align"`
	fieldRow    *js.Object `js:"fieldRow"`     // []*Blockly.Field
	sourceBlock *js.Object `js:"sourceBlock_"` // *Blockly.Block
	connection  *js.Object `js:"connection"`   // *Blockly.Connection
}

func (in *Input) Block() *Block {
	return &Block{Object: in.sourceBlock}
}

func (in *Input) Connection() *Connection {
	return jsConnection(in.connection)
}

// blockly's append field allows field to be a string, and then to pass an optional name
// see also: insertFieldAt
func (in *Input) AppendField(f *Field) {
	in.Call("appendField", f.Object)
}

func (in *Input) AppendNamedField(name string, f *Field) {
	in.Call("appendField", f.Object, name)
}

func (in *Input) Fields() (ret *Fields) {
	if obj := in.fieldRow; obj != nil && obj.Bool() {
		ret = &Fields{Object: obj}
	}
	return
}

// insertFieldAt = function(index, field, OptName) {
// removeField = function(name) {

func (in *Input) IsVisible() bool {
	obj := in.Call("isVisible")
	return obj.Bool()
}

func (in *Input) SetVisible(visible bool) {
	in.Call("setVisible", visible)
}

var invisible = js.MakeFunc(func(*js.Object, []*js.Object) (ret interface{}) {
	return false
})

func (in *Input) SetCheck(compatibleType block.Type) (err error) {
	var ar []block.Type
	if compatibleType != "" {
		ar = append(ar, compatibleType)
	}
	return in.SetChecks(ar)
}

func (in *Input) SetChecks(compatibleTypes []block.Type) (err error) {
	// pattern for handling thrown errors
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(*js.Error); ok {
				err = errutil.New(e, in.Name, in.Type, compatibleTypes)
			} else {
				panic(e)
			}
		}
	}()
	in.Call("setCheck", compatibleTypes)
	return
}

func (in *Input) SetAlign(a InputAlign) {
	in.Call("setAlign", a)
}

func (in *Input) Dispose() {
	in.Call("dispose")
}
