package gblocks

import (
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
	"github.com/ionous/gblocks/test"
	"github.com/ionous/gblocks/tin"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

// Build a description of test.InputBlock.
// InputBlock contains the three blockly input types: value, statement, and dummy.
// ( package mutant uses input dummy(s) for generic mutation data. )
func TestInputs(t *testing.T) {
	in := (*test.InputBlock)(nil)
	row := (*test.RowBlock)(nil)
	mutation := (*test.BlockMutation)(nil)
	ms := tin.Mutables{}
	if e := ms.AddMutation(mutation); e != nil {
		t.Fatal(e)
	} else if mutation, ok := ms.FindMutable("block_mutation"); !ok {
		t.Fatal("cant find mutation")
	} else {
		if types, e := new(TypeCollector).
			AddTerm(row).
			AddStatement(in).
			GetTypes(); e != nil {
			t.Fatal(e)
		} else {
			m := Maker{types: types, mutables: ms}
			var mins mutant.InMutations
			if desc, e := m.makeDesc("input_block", &mins); e != nil {
				t.Fatal(e)
			} else {
				expectedDesc := block.Dict{
					"type":              "input_block",
					"tooltip":           "input block",
					"previousStatement": "input_block",
					"message0":          "%1 %2 %3",
					"args0": []block.Dict{
						{
							"name":  "VALUE",
							"type":  "input_value",
							"check": "row_block",
						},
						{
							"name":  "STATEMENT",
							"type":  "input_statement",
							"check": "input_block",
						},
						{
							"name": "MUTATION",
							"type": "input_dummy",
						},
					},
				}
				if v := pretty.Diff(expectedDesc, desc); len(v) != 0 {
					t.Log(pretty.Print(desc))
					t.Fatal(v)
				}
				require.Equal(t, mins.Inputs, []string{"MUTATION"})
				require.EqualValues(t, mins.Mutations["MUTATION"], mutation)
				require.Equal(t, "block_mutation", mutation.Name())
			}
		}
	}
}
