package block

import (
	r "reflect"
	"strconv"
)

// Item - used for inputs and fields. assumes caps case. ex. INPUT_NAME
type Item string

func ItemFromField(f r.StructField) Item {
	name := pascalToCaps(f.Name)
	return ItemFromString(name)
}

func ItemFromString(s string) Item {
	return Item(s)
}

func (n Item) Push(o Item) (ret Item) {
	if len(n) > 0 {
		ret = n + "/" + o
	} else {
		ret = o
	}
	return
}

func (n Item) Index(i int) Item {
	s := strconv.Itoa(i)
	return n.Push(Item(s))
}

// Friendly returns the name in spaces. ex. "Item Name"
func (n Item) Friendly() string {
	return pascalToSpace(underscoreToPascal(string(n)))
}

// String returns the name in default (caps ) ex. "INPUT_NAME"
func (n Item) String() string {
	return string(n)
}
