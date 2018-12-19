package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type InputType int

//go:generate stringer -type=InputType
const (
	InputValue InputType = iota + 1
	OutputValue
	NextStatement
	PreviousStatement
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
	*js.Object               // Blockly.Input
	Type       InputType     `js:"type"`
	Name       string        `js:"name"`
	Align      InputAlign    `js:"align"`
	FieldRow   []interface{} `js:"fieldRow"`   // array of Blockly.Field
	Connection *Connection   `js:"connection"` // Blockly.Connection
	// custom field to count the number of following inputs resulting from a mutation
	// 0 ( the default ) means, this input will never have mutations;
	// -1 means it can have mutations, but currently does not
	mutations int `js:"mutations_"`
}

// blockly's append field allows field to be a string, and then to pass an optional name
// see also: insertFieldAt
func (in *Input) AppendField(f *Field) {
	in.Call("appendField", f.Object)
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

func (in *Input) ForceMutation() {
	in.Set("isVisible", invisible)
	in.SetVisible(false)
	in.mutations = -1
}

func (in *Input) SetCheck(compatibleType string) (err error) {
	var ar []string
	if compatibleType != "" {
		ar = []string{compatibleType}
	}
	return in.SetChecks(ar)
}

func (in *Input) SetChecks(compatibleTypes []string) (err error) {
	// handle thrown errors for missing connections, etc.
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(*js.Error); ok {
				err = e
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

// init = function() {

func (in *Input) Dispose() {
	in.Call("dispose")
}
