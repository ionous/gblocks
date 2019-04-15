// extremely minimial wrapper for javascript dom
package jsdom

import (
	"github.com/gopherjs/gopherjs/js"
)

type Element struct {
	*js.Object
}

func (m *Element) OuterHTML() string {
	return m.Get("outerHTML").String()
}
