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

func (cs *Connections) AppendInput(in *Input) {
	var target *Connection
	if c := in.Connection(); c != nil {
		target = c.TargetConnection()
	}
	cs.SetIndex(cs.Object.Length(), target)
}

func (cs *Connections) Connection(i int) *Connection {
	return jsConnection(cs.Index(i))
}
