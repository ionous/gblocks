package gblocks

import (
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/test"
	"github.com/ionous/gblocks/tin"
	"github.com/stretchr/testify/require"
)

// expand an atom that contains an interface
func TestInterface(t *testing.T) {
	// 1. register a mutation that has an input with an interface
	ms := tin.Mutables{}
	if e := ms.AddMutation((*test.BlockMutation)(nil),
		(*test.AtomWithInterface)(nil),
	); e != nil {
		t.Fatal("add the test mutation", e)
	} else if types, e := new(TypeCollector).
		AddTerm((*test.InterfacingTerm)(nil)).
		AddTerm((*test.MutableBlock)(nil)).
		GetTypes(); e != nil {
		t.Fatal(t, e)
		// find the just added mutation
	} else if mutable, ok := ms.FindMutable("block_mutation"); !ok {
		t.Fatal("find the mutation")
		// find the "atom with interface" quark
	} else if q, ok := mutable.Quarks(true); !ok {
		t.Fatal("get the quark")
		// expand the "atom with interfac"
	} else if args, e := q.Atomize("scoped", &Maker{types: types}); e != nil {
		t.Fatal("expand atom")
	} else {
		// we expect "atom with interface" has one input;
		// of the two terms we added -- 'mutable block' and 'interfacing term' --
		// only the latter should be attachable to its input.
		expected := block.NewArgs("%1", block.Dict{
			"name":  block.Scope("scoped", "INPUT"),
			"type":  "input_value",
			"check": "interfacing_term",
		})
		require.Equal(t, expected, args)
	}
}
