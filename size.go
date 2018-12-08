package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

// goog.math.Size
type Size struct {
	*js.Object
	Width  string `js:"width"`
	Height string `js:"width"`
}
