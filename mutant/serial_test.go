package mutant_test

import (
	"testing"

	"github.com/ionous/gblocks/dom"
	"github.com/ionous/gblocks/mock"
	"github.com/ionous/gblocks/mutant"
	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) {
	data := dom.BlockMutation{dom.Mutations{
		&dom.Mutation{
			Input: "I1",
			Atoms: dom.Atoms{
				[]string{"a1"},
			},
		},
		&dom.Mutation{
			Input: "I3",
			Atoms: dom.Atoms{
				[]string{"a2", "a3"},
			},
		},
	}}
	expectedInputs := mutant.MutableInputs{
		"I1": {"a1"},
		"I3": {"a2", "a3"},
	}
	expandedInputs := []string{
		"I1:input_dummy",
		/**/ "a$mock$I1$0$TERM:input_value",
		"I2:input_dummy",
		"I3:input_dummy",
		/**/ "a$mock$I3$0$NUM:field_number",
		/**/ "a$mock$I3$1$TEXT:field_input",
		/**/ "a$mock$I3$1$STATE:input_statement",
	}
	expectedOutput := `` +
		`<mutation>` +
		/**/ `<pin name="I1">` +
		/* */ `<atom type="a1"></atom>` +
		/**/ `</pin>` +
		/**/ `<pin name="I3">` +
		/* */ `<atom type="a2"></atom>` +
		/* */ `<atom type="a3"></atom>` +
		/**/ `</pin>` +
		`</mutation>`

	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
	mutator := mock.NewInMutations(common.inputs, common.quarks)

	db := &mock.MockDatabase{common.atomProducts}
	inputs, e := mutator.LoadMutation(b, db, data)
	require.NoError(t, e)
	require.Equal(t, expectedInputs, inputs)
	expanded := listInputs(b)
	require.Equal(t, expandedInputs, expanded)

	// save the in-memory data, makes sure it matches the original input
	serial := mutator.SaveMutation(inputs)
	require.Equal(t, data, serial)
	if str, e := serial.MarshalMutation(); e != nil {
		t.Fatal(e)
	} else {
		require.Equal(t, expectedOutput, str)
	}
}
