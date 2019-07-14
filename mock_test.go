package gblocks

import (
	"testing"

	"github.com/ionous/gblocks/mock"
	"github.com/ionous/gblocks/mutant"
	"github.com/ionous/gblocks/test"
	"github.com/ionous/gblocks/tin"
	"github.com/stretchr/testify/require"
)

// create a mocked up block ( rather than blockly which requires javascript )
func TestMockBlock(t *testing.T) {
	in := (*test.InputBlock)(nil)
	row := (*test.RowBlock)(nil)
	mutation := (*test.BlockMutation)(nil)
	ms := tin.Mutations{}
	if e := ms.AddMutation(mutation); e != nil {
		t.Fatal(e)
	}
	if types, e := new(TypeCollector).
		AddTerm(row).
		AddStatement(in).
		GetTypes(); e != nil {
		t.Fatal(e)
	} else {
		m := Maker{types: types, mutables: ms}
		var ins mutant.BlockMutations
		if desc, e := m.makeDesc("input_block", &ins); e != nil {
			t.Fatal(e)
		} else {
			b := mock.CreateBlock("mock", desc)
			expected := []struct{ Name, Type string }{{
				Name: "VALUE",
				Type: "input_value",
			}, {
				Name: "STATEMENT",
				Type: "input_statement",
			}, {
				Name: "MUTATION",
				Type: "input_dummy",
			},
			}
			require.Equal(t, b.BlockType(), "input_block")
			require.Equal(t, 3, b.NumInputs())
			for i, x := range expected {
				in := b.Input(i)
				require.Equal(t, in.InputName(), x.Name)
				require.Equal(t, in.InputType(), x.Type)
			}
		}
	}
}
