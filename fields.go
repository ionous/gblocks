package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type FieldAngle float32
type FieldCheckbox bool
type FieldColour string //'#rrggbb'
//etype FieldDate string   //  goog.date.Date().toIsoString(true);

// options seem to be image or text
type FieldDropdown []string
type FieldImageDropdown []FieldImage

type FieldImage struct {
	Width, Height int
	Src           string
	Alt           string
}

type FieldText string // field_input, FieldTextInput; pre-existing validators inclde numberValidator, nongenativeIntegerValidator
type FieldNumber float32
type FieldVariable string

type Fields struct {
	*js.Object
}

func (f *Fields) Field(i int) (ret *Field) {
	if obj := f.Index(i); obj != nil && obj.Bool() {
		ret = &Field{Object: obj}
	}
	return
}
