package mock

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

type Registry struct {
	Blocks   map[string]block.Dict
	Mutators map[string]block.Mutator
}

func (reg *Registry) IsBlockRegistered(blockType string) bool {
	_, ok := reg.Blocks[blockType]
	return ok
}

func (reg *Registry) RegisterBlock(blockType string, desc block.Dict) (err error) {
	if _, ok := reg.Blocks[blockType]; ok {
		err = errutil.New(blockType, "already registered")
	} else if reg.Blocks == nil {
		reg.Blocks = map[string]block.Dict{blockType: desc}
	} else {
		reg.Blocks[blockType] = desc
	}
	return
}

func (reg *Registry) NewMockSpace() *MockSpace {
	return &MockSpace{reg, make(map[string]block.Shape), make(map[string]int), nil}
}

func (reg *Registry) RegisterMutator(name string, m block.Mutator) (err error) {
	if _, ok := reg.Mutators[name]; ok {
		err = errutil.New(name, "already registered")
	} else if reg.Mutators == nil {
		reg.Mutators = map[string]block.Mutator{name: m}
	} else {
		reg.Mutators[name] = m
	}
	return
}
