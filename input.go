package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	"strings"
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

// InputNames are caps case. ex. INPUT_NAME
type InputName string

// FieldPath returns the name of the go struct field.
func (n InputName) FieldPath() string {
	slashes := strings.Split(n.String(), "/")
	for i, s := range slashes {
		slashes[i] = underscoreToPascal(s)
	}
	return strings.Join(slashes, "/")
}

// Friendly returns the name in spaces.
func (n InputName) Friendly() string {
	return pascalToSpace(underscoreToPascal(n.String()))
}

// String returns the name in default (uppercase)
func (n InputName) String() (ret string) {
	if len(n) > 0 {
		ret = string(n)
	}
	return
}

type Input struct {
	*js.Object               // Blockly.Input
	Type       InputType     `js:"type"`
	Name       InputName     `js:"name"`
	Align      InputAlign    `js:"align"`
	FieldRow   []interface{} `js:"fieldRow"`   // array of Blockly.Field
	connection *js.Object    `js:"connection"` // Blockly.Connection
	mutation_  *js.Object    `js:"mutation_"`  // *InputMutation
}

func (in *Input) Connection() *Connection {
	return jsConnection(in.connection)
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

func (in *Input) ForceMutation(name string) {
	in.Set("isVisible", invisible)
	in.SetVisible(false)
	m := &InputMutation{Object: new(js.Object)}
	m.name = name
	in.mutation_ = m.Object
}

func (in *Input) Mutation() (ret *InputMutation) {
	if obj := in.mutation_; obj != nil && obj.Bool() {
		ret = &InputMutation{Object: in.mutation_}
	}
	return
}

func (in *Input) SetCheck(compatibleType TypeName) (err error) {
	var ar []TypeName
	if compatibleType != "" {
		ar = append(ar, compatibleType)
	}
	return in.SetChecks(ar)
}

func (in *Input) SetChecks(compatibleTypes []TypeName) (err error) {
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

// init = function() {

func (in *Input) Dispose() {
	in.Call("dispose")
}
