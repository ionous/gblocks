package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type Connections struct {
	*js.Object
}

func NewConnections() *Connections {
	return &Connections{Object: js.Global.Get("Array").New()}
}

func (cs *Connections) Append(c *Connection) {
	cs.SetIndex(cs.Length(), c)
}

func (cs *Connections) Connection(i int) (ret *Connection) {
	return jsConnection(cs.Index(i))
}

// it's not clear why, but using (an uninitialied) *Connections as a js:tagged field
// results in valid *Connections pointer with a nil *js.Object.
func (cs *Connections) Length() (ret int) {
	if cs.Object != nil && cs.Object.Bool() {
		ret = cs.Object.Length()
	}
	return
}
