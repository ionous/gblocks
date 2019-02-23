package option

import "strconv"

const (
	Output    = "output"
	Class     = "class"
	Prev      = "previousStatement"
	Next      = "nextStatement"
	Name      = "name"
	Type      = "type"
	Check     = "check"
	Choices   = "options"
	Value     = "value"
	Min       = "min"
	Max       = "max"
	Text      = "text"
	Precision = "precision"
	// custom options
	Group = "mutation"
	Decor = "decor"
)

func Message(i int) string {
	return "message" + strconv.Itoa(i)
}

func Args(i int) string {
	return "args" + strconv.Itoa(i)
}
