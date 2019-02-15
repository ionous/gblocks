package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type Field struct {
	*js.Object
	name *js.Object `js:"name"` // note: added on by appendField; so sometimes undefined.
	//maxDisplayLength
}

type FieldLabel struct {
	*Field
	class string `js:"class_"`
}

// opt_class Optional CSS class for the field's text.
func NewFieldLabel(txt, cls string) (ret *FieldLabel) {
	if blockly := GetBlockly(); blockly != nil {
		obj := blockly.Get("FieldLabel").New(txt, cls)
		ret = &FieldLabel{Field: &Field{Object: obj}}
	}
	return
}

// note: added, optionally, on by appendField; sometimes undefined.
func (f *Field) Name() (ret string) {
	if name := f.name; name != nil && name.Bool() {
		ret = name.String()
	}
	return
}

// func (f *Field) setSourceBlock (block) {
// func (f *Field) init () {
// func (f *Field) initModel () {
// func (f *Field) dispose () {
// func (f *Field) updateEditable () {
// func (f *Field) isCurrentlyEditable () {
// func (f *Field) isVisible () {
// func (f *Field) setVisible (visible) {
// func (f *Field) setValidator (handler) {
// func (f *Field) getValidator () {
// func (f *Field) classValidator (text) {
// func (f *Field) callValidator (text) {
// func (f *Field) getSvgRoot () {
// func (f *Field) updateWidth () {
func (f *Field) GetSize() *Size {
	obj := f.Call("getSize")
	return &Size{Object: obj}
}

func (f *Field) GetText() (ret string) {
	if obj := f.Call("getText"); obj.Bool() {
		ret = obj.String()
	}
	return
}

func (f *Field) SetText(newText string) {
	f.Call("setText", newText)
}

// func (f *Field) forceRerender () {
func (f *Field) GetValue() *js.Object {
	return f.Call("getValue")
}

// in blockly, by default, this routes to setText
func (f *Field) SetValue(newValue interface{}) {
	f.Call("setValue", newValue)
}

// in blockly newTip can be an element as well
func (f *Field) SetTooltip(newTip string) {
	f.Call("setTooltip", newTip)
}

// func (f *Field) referencesVariables () {
