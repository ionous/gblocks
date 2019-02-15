package named

import r "reflect"

// Item - used for inputs and fields. assumes caps case. ex. INPUT_NAME
type Item string

func ItemFromField(f r.StructField) Item {
	name := pascalToCaps(f.Name)
	return Item(name)
}

// Friendly returns the name in spaces. ex. "Item Name"
func (n Item) Friendly() string {
	return pascalToSpace(underscoreToPascal(n.String()))
}

// String returns the name in default (caps ) ex. "INPUT_NAME"
func (n Item) String() (ret string) {
	if len(n) > 0 {
		ret = string(n)
	}
	return
}
