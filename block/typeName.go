package block

import (
	r "reflect"
	"strings"
)

type Type string

// SpecialType - generate a name from separated parts (using separators which use javascript friendly, golang impossible identifiers)
func SpecialType(parts ...string) Type {
	name := strings.Join(parts, "$")
	return Type(name)
}

func TypeFromStruct(t r.Type) Type {
	return TypeFromStructName(t.Name())
}

// PascalCase -> under_score
func TypeFromStructName(structName string) Type {
	name := pascalToUnderscore(structName)
	return Type(name)
}

// StructName - returns a go-style name.
// ex. "ExampleBlock"
func (n Type) StructName() string {
	return underscoreToPascal(n.String())
}

// StructName - returns a human readable name.
// ex. "example block"
func (n Type) Friendly() string {
	return pascalToSpace(underscoreToPascal(n.String()))
}

// String - returns a blockly-style type name.
// ex. "example-block"
func (n Type) String() (ret string) {
	if len(n) > 0 {
		ret = string(n)
	}
	return
}
