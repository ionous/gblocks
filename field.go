package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type Field struct {
	*js.Object
	Name string `js:"name"`
}
