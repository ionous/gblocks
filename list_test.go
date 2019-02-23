package gblocks

import (
	// "github.com/kr/pretty"
	// "github.com/stretchr/testify/require"
	// r "reflect"
	"testing"
)

type ListTest struct {
	ListMutation
}

type NeverEmptyListTest struct {
	NeverEmptyListMutation
}

type ListMutation struct {
	NextStatement *ListElement
}

type NeverEmptyListMutation struct {
	Element       interface{} // any block
	NextStatement *ListElement
}

type ListElement struct {
	Element interface{}
	// PrevStatement
	NextStatement *ListElement
}

func TestListMutations(t *testing.T) {
	// var reg Registry
	// require.NoError(t,
	// 	reg.RegisterMutation("ListMutation",
	// 		Mutation{"item", (*ListElement)(nil)},
	// 	))
	// require.NoError(t,
	// 	reg.RegisterMutation("NeverEmptyListMutation",
	// 		Mutation{"item", (*ListElement)(nil)},
	// 	))
	// //
	// require.NoError(t, reg.RegisterBlocks(nil,
	// 	(*ListTest)(nil),
	// 	(*NeverEmptyListTest)(nil),
	// 	(*ListElement)(nil),
	// ), "register blocks")
	// ws := NewBlankWorkspace(false,&orderedGenerator{name: "lists"})
	// // for testing, replace timed event queue with direct event queue
	// blockly.Events().EnableTestFiring()
	// // fn(ws, reg)
	// ws.Dispose()
}
