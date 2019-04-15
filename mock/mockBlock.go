package mock

import (
	"fmt"

	"github.com/ionous/gblocks/block"
)

type MockBlock struct {
	Id   string
	Type string
	InputList
	Next, Prev *MockConnection
	Workspace  *MockSpace
	Flags      Flags
}

type Flags map[block.Flag]bool

func (b *MockBlock) BlockId() string   { return b.Id }
func (b *MockBlock) BlockType() string { return b.Type }
func (b *MockBlock) SetFlag(f block.Flag, v bool) {
	b.Flags[f] = v
}

func (b *MockBlock) GetFlag(f block.Flag) bool {
	return b.Flags[f]
}

func (b *MockBlock) String() string {
	return fmt.Sprintf("%s:%s", b.Id, b.Type)
}

func (b *MockBlock) HasWorkspace() bool { return true }
func (b *MockBlock) BlockWorkspace() (ret block.Workspace) {
	if b.Workspace != nil {
		ret = b.Workspace
	}
	return
}
func (b *MockBlock) InitSvg() {}

func (b *MockBlock) Dispose() {
	if ws := b.Workspace; ws != nil {
		ws.Delete(b.Id)
		b.Workspace = nil
	}
}

// connection to a piece in the following line
func (b *MockBlock) NextConnection() (ret block.Connection) {
	if b.Next != nil {
		ret = b.Next
	}
	return
}

// connection to a piece in the following line
func (b *MockBlock) PreviousConnection() (ret block.Connection) {
	if b.Prev != nil {
		ret = b.Prev
	}
	return
}
