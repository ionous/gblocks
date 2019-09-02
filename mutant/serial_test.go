package mutant_test

import (
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
	"github.com/ionous/gblocks/mock"
	"github.com/ionous/gblocks/mutant"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) {
	data := dom.BlockMutation{dom.Mutations{
		&dom.Mutation{
			Input: "M1",
			Atoms: dom.Atoms{
				// A -> Term:value_input
				&dom.Atom{"name1", "A"},
			},
		},
		&dom.Mutation{
			Input: "M3",
			Atoms: dom.Atoms{
				// B -> Num:number_field
				&dom.Atom{"name2", "B"},
				// C -> Text:text_field, State:statement_input
				&dom.Atom{"name3", "C"},
			},
		},
	}}
	expandedInputs := []string{
		"M1:input_dummy",
		/**/ block.Scope("a", "name1", "TERM:input_value"),
		"M2:input_dummy",
		"M3:input_dummy",
		/**/ block.Scope("a", "name2", "NUM:field_number"),
		/**/ block.Scope("a", "name3", "TEXT:field_input"),
		/**/ block.Scope("a", "name3", "STATE:input_statement"),
	}
	expectedOutput := `` +
		`<mutation>` +
		/**/ `<pin name="M1">` +
		/* */ `<atom name="name1" type="A"></atom>` +
		/**/ `</pin>` +
		/**/ `<pin name="M3">` +
		/* */ `<atom name="name2" type="B"></atom>` +
		/* */ `<atom name="name3" type="C"></atom>` +
		/**/ `</pin>` +
		`</mutation>`

	// create a block with three mutable inputs: M1, M2, M3
	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
	arch := mock.NewMutations(common.inputs, common.quarks)
	// "database" to create workspace inputs from types of atoms
	db := &mock.MockDatabase{common.atomProducts}
	blocks := mutant.NewMutatedBlocks()
	// create a mutated block record for our mock block
	mb := blocks.CreateMutatedBlock(b, arch, db)
	e := mb.LoadMutation(&data)
	require.NoError(t, e)
	expanded := listInputs(b)
	require.Equal(t, expandedInputs, expanded)

	// save the in-memory data, makes sure it matches the original input
	serial := mb.SaveMutation()
	if v := pretty.Diff(data, serial); len(v) != 0 {
		t.Log(pretty.Sprint(serial))
		t.Fatal(v)
	}
	if str, e := serial.MarshalMutation(); e != nil {
		t.Fatal(e)
	} else {
		require.Equal(t, expectedOutput, str)
	}
}
