package option

import "strconv"

const (
	Check     = "check"
	Choices   = "options"
	Class     = "class"
	Max       = "max"
	Min       = "min"
	Mutator   = "mutator"
	Name      = "name"
	Next      = "nextStatement"
	Output    = "output"
	Precision = "precision"
	Prev      = "previousStatement"
	Text      = "text"
	Type      = "type"
	Value     = "value"
	Input     = "input"
	Colour    = "colour"
	Tooltip   = "tooltip"
	Help      = "helpUrl"
	// custom options
	Decor = "decor"
)

func Message(i int) string {
	return "message" + strconv.Itoa(i)
}

func Args(i int) string {
	return "args" + strconv.Itoa(i)
}
