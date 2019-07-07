package option

import (
	"strconv"
	"strings"
)

const (
	Check        = "check"
	Choices      = "options"
	Class        = "class"
	Colour       = "colour"
	Help         = "helpUrl"
	Input        = "input"
	InputsInline = "inputsInline"
	Max          = "max"
	Min          = "min"
	Mutator      = "mutator"
	Name         = "name"
	Next         = "nextStatement"
	Output       = "output"
	Precision    = "precision"
	Prev         = "previousStatement"
	Text         = "text"
	Tooltip      = "tooltip"
	Type         = "type"
	Value        = "value"
	// custom options
	Decor = "decor"
)

func Message(i int) string {
	return "message" + strconv.Itoa(i)
}

func Args(i int) string {
	return "args" + strconv.Itoa(i)
}

// extract the desired blockly input type ( a string name ) from a dictionary of struct tags.
// w/o a tag, we wind up with an "input value" ( a generic term )
// fix? if really needed to decouple this, could have an "input" factory
func InputOption(src interface {
	Lookup(key string) (ret string, okay bool)
}) (ret string, okay bool) {
	tag, ok := src.Lookup(Input)
	lower := strings.ToLower(tag)
	switch lower {
	case "":
		lower = "value"
	case "mutation", "option", "choice", "repetition":
		lower = "dummy"
	}
	out := "input_" + lower
	return out, ok
}
