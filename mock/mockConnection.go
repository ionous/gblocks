package mock

import "github.com/ionous/gblocks/block"

type MockConnection struct {
	Name   string
	Source block.Shape
	Target block.Connection
}

func (c *MockConnection) String() (ret string) {
	ret = c.Name + ":"
	if tgt := c.TargetConnection(); tgt != nil {
		ret += tgt.(*MockConnection).Name
	}
	return
}

func (c *MockConnection) SourceBlock() block.Shape {
	return c.Source
}

//  Connection this connection connects to.  Null if not connected.
func (c *MockConnection) TargetConnection() block.Connection {
	return c.Target
}

func (c *MockConnection) TargetBlock() (ret block.Shape) {
	if tc := c.Target; tc != nil {
		ret = tc.SourceBlock()
	}
	return
}

func (c *MockConnection) IsConnected() bool {
	return c.Target != nil
}

// blockly does *a lot* of work around shadow blocks,
// bumping connections to "inject" the new connection "in front" of existing connections
// orphaning blocks, recording undos, events, etc.
func (c *MockConnection) Connect(other block.Connection) {
	if other == nil {
		panic("Cant connect to nothing")
	} else if other != c.Target {
		if other.IsConnected() {
			other.Disconnect()
		}
		if c.IsConnected() {
			c.Disconnect()
		}
		other := other.(*MockConnection)
		c.Target, other.Target = other, c
	}
}

// in blockly determines which is "superior" (aka. parent)
// [ higher in a stack or receiving input rather than generating output ]
// then fires BlockMove parent to child
func (c *MockConnection) Disconnect() {
	if other := c.Target.(*MockConnection); other == nil {
		panic("Source connection not connected.")
	} else if other.Target != c {
		panic("Target connection not connected to source connection.")
	} else {
		c.Target, other.Target = nil, nil
	}
}
