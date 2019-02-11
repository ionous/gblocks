package named

import r "reflect"

// Input - assumes caps case. ex. INPUT_NAME
type Input string

func InputFromField(f r.StructField) Input {
	name := pascalToCaps(f.Name)
	return Input(name)
}

// Friendly returns the name in spaces. ex. "Input Name"
func (n Input) Friendly() string {
	return pascalToSpace(underscoreToPascal(n.String()))
}

// String returns the name in default (caps ) ex. "INPUT_NAME"
func (n Input) String() (ret string) {
	if len(n) > 0 {
		ret = string(n)
	}
	return
}
