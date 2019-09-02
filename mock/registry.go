package mock

import (
	"fmt"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

type MockProject struct {
	Blocks    map[string]block.Dict
	Mutators  map[string]block.Mutator
	Uniquer   UniqueId
	Workspace UniqueId
}

type UniqueId struct{ id int }

func (u *UniqueId) GenerateUniqueName(context string) string {
	n := u.id
	u.id++
	return fmt.Sprintf("%s-%05d", context, n)
}

func (reg *MockProject) GenerateUniqueName() string {
	return reg.Uniquer.GenerateUniqueName("name")
}

func (reg *MockProject) IsBlockRegistered(blockType string) bool {
	_, ok := reg.Blocks[blockType]
	return ok
}

func (reg *MockProject) RegisterBlock(blockType string, desc block.Dict) (err error) {
	if _, ok := reg.Blocks[blockType]; ok {
		err = errutil.New(blockType, "already registered")
	} else if reg.Blocks == nil {
		reg.Blocks = map[string]block.Dict{blockType: desc}
	} else {
		reg.Blocks[blockType] = desc
	}
	return
}

func (reg *MockProject) NewMockSpace() *MockSpace {
	wsid := reg.Workspace.GenerateUniqueName("ws")
	return &MockSpace{reg.Blocks, wsid, make(map[string]block.Shape), make(map[string]int), nil}
}

func (reg *MockProject) RegisterMutator(name string, m block.Mutator) (err error) {
	if _, ok := reg.Mutators[name]; ok {
		err = errutil.New(name, "already registered")
	} else if reg.Mutators == nil {
		reg.Mutators = map[string]block.Mutator{name: m}
	} else {
		reg.Mutators[name] = m
	}
	return
}
