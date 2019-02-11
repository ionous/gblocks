package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

// Connection *potentially* connects to another block; it's more a "connector" than an "connection".
type Connection struct {
	*js.Object
	Type             InputType  `js:"type"`
	targetConnection *js.Object `js:"targetConnection"`
}

func (c *Connection) GetSourceBlock() (ret *Block) {
	if obj := c.Call("getSourceBlock"); obj != nil {
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

func (c *Connection) Connect(o *Connection) {
	c.Call("connect", o.Object)
}

func (c *Connection) Disconnect() {
	c.Call("disconnect")
}

func (c *Connection) TargetBlock() (ret *Block) {
	if obj := c.Call("targetBlock"); obj != nil {
		ret = &Block{Object: obj}
	}
	return
}
func (c *Connection) TargetConnection() *Connection {
	return jsConnection(c.targetConnection)
}

//func (c*Connection) setCheck(c[]string) {// array or value
// func (c*Connection) getCheck() *js.Object {// array
// func (c*Connection) setShadowDom(el)
// func (c*Connection) getShadowDom()// el
// toString() string
