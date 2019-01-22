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
	cs.SetIndex(c.Length(), c.Object)
}

func (cs *Connections) Connection(i int) *Connection {
	return &Connection{Object: cs.Index(i)}
}
