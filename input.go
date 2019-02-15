package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/named"
)

// InputType - describes both inputs and connections
type InputType int

//go:generate stringer -type=InputType
const (
	InputValue InputType = iota + 1
	OutputValue
	NextStatement     // used for connections between blocks, and for statement inputs
	PreviousStatement // used for connections between blocks
	DummyInput
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
	Name        named.Item `js:"name"`
	Align       InputAlign `js:"align"`
	fieldRow    *js.Object `js:"fieldRow"`     // []*Blockly.Field
	sourceBlock *js.Object `js:"sourceBlock_"` // *Blockly.Block
	connection  *js.Object `js:"connection"`   // *Blockly.Connection
	// custom
	mutation_ *js.Object `js:"mutation_"` // *InputMutation
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

// insertFieldAt = function(index, field, opt_name) {
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

// ForceMutation - name is mutation name. see RegisterMutation
func (in *Input) ForceMutation(name named.Type) *InputMutation {
	in.Set("isVisible", invisible)
	in.SetVisible(false)
	mutation := NewInputMutation(in, name)
	in.mutation_ = mutation.Object
	return mutation
}

func (in *Input) Mutation() (ret *InputMutation) {
	if obj := in.mutation_; obj != nil && obj.Bool() {
		ret = &InputMutation{Object: in.mutation_}
	}
	return
}

func (in *Input) SetCheck(compatibleType named.Type) (err error) {
	var ar []named.Type
	if compatibleType != "" {
		ar = append(ar, compatibleType)
	}
	return in.SetChecks(ar)
}

func (in *Input) SetChecks(compatibleTypes []named.Type) (err error) {
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

// iterate over all blocks stacked in this input
func (in *Input) visitStack(cb func(b *Block) (keepGoing bool)) (exhausted bool) {
	earlyOut := false
	// get the input's connection information
	if c := in.Connection(); c != nil {
		// for every block connected to the input...
		for b := c.TargetBlock(); b != nil; {
			if !cb(b) {
				earlyOut = true
				break
			}

			// move to the next
			if c := b.NextConnection(); c != nil {
				b = c.TargetBlock()
			} else {
				break
			}
		}
	}
	return !earlyOut
}
