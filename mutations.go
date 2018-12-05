package gblocks

import (
	r "reflect"
)

// Mutation - base class for blocks used in mutation dialogs.
type Mutation struct {
	BlockType  r.Type
	Connection *Connection
}

type Mutator interface {
	// blocks which appear in the mutation dialog's palette
	Types() []r.Type
	NumMutations() int
	Mutations() []Mutation // possibly r.Type
	NewMutations() []Mutation
}

func NewMutation(i interface{}) *Mutation {
	return &Mutation{
		BlockType: r.TypeOf(i),
	}
}
