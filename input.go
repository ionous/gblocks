package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type InputType int

const (
	InputValue InputType = iota + 1
	OutputValue
	NextStatement
	PreviousStatement
	DummyInput
)

type InputAlign int

const (
	AlignLeft InputAlign = iota - 1
	AlignCentre
	AlignRight
)

type Input struct {
	*js.Object               // Blockly.Input
	Type       int           `js:"type"`
	Name       string        `js:"string"`
	Align      InputAlign    `js:"align"`
	FieldRow   []interface{} `js:"fieldRow"`   // array of Blockly.Field
	Connection *Connection   `js:"connection"` // Blockly.Connection
}

// blockly's append field allows field to be a string, and then to pass an optional name
// see also: insertFieldAt
func (i *Input) AppendField(f *Field) (err error) {
	i.Call("appendField", f.Object)
	return
}

func (i *Input) SetAlign(a InputAlign) (err error) {
	i.Call("setAlign", a)
	return
}

func (i *Input) SetCheck(compatibleType string) (err error) {
	var ar []string
	if compatibleType != "" {
		ar = []string{compatibleType}
	}
	return i.SetChecks(ar)
}

func (i *Input) SetChecks(compatibleTypes []string) (err error) {
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
	i.Call("setCheck", compatibleTypes)
	return
}
