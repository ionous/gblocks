package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/block"
)

// Connection *potentially* connects to another block; it's more a "connector" than an "connection".
type Connection struct {
	*js.Object
	Type             InputType  `js:"type"`
	targetConnection *js.Object `js:"targetConnection"`
}

func (c *Connection) SourceBlock() (ret block.Shape) {
	if obj := c.Call("getSourceBlock"); obj.Bool() {
		ret = &Block{Object: obj}
	}
	return
}

func (c *Connection) IsSuperior() bool {
	return c.Call("isSuperior").Bool()
}

func (c *Connection) IsConnected() bool {
	return c.Call("isConnected").Bool()
}

func (c *Connection) IsConnectionAllowed() bool {
	return c.Call("isConnectionAllowed").Bool()
}

func (c *Connection) Connect(o block.Connection) {
	c.Call("connect", o.(*Connection).Object)
}

func (c *Connection) Disconnect() {
	c.Call("disconnect")
}

func (c *Connection) TargetBlock() (ret block.Shape) {
	if obj := c.Call("targetBlock"); obj.Bool() {
		ret = &Block{Object: obj}
	}
	return
}
func (c *Connection) TargetConnection() block.Connection {
	return jsConnection(c.targetConnection)
}

//func (c*Connection) setCheck(c[]string) {// array or value
// func (c*Connection) getCheck() *js.Object {// array
// func (c*Connection) setShadowDom(el)
// func (c*Connection) getShadowDom()// el
// toString() string
