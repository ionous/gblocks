package gblocks

import (
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/enum"
	"github.com/ionous/gblocks/test"
	"github.com/stretchr/testify/require"
)

func TestDescWithEnum(t *testing.T) {
	// make from EnumStatement
	var pairs enum.Pairs
	_, e := pairs.AddEnum(map[test.Enum]string{
		test.DefaultChoice:     "default",
		test.AlternativeChoice: "alt",
	})
	require.NoError(t, e)
	if types, e := new(TypeCollector).
		AddTopStatement((*test.EnumStatement)(nil)).
		GetTypes(); e != nil {
		t.Fatal(e)
	} else {
		m := Maker{pairs: pairs, types: types}
		desc, e := m.makeDesc("enum_statement", nil)
		require.NoError(t, e)
		expected := block.Dict{
			"message0": "%1",
			"type":     "enum_statement",
			"args0": []block.Dict{
				{
					"name": "ENUM",
					"type": "field_dropdown",
					"options": []enum.Pair{
						{"default", "DefaultChoice"},
						{"alt", "AlternativeChoice"},
					},
				},
			},
		}
		require.Equal(t, expected, desc)
	}
}

func TestDescWithOutput(t *testing.T) {
	var types TypeCollector
	types.AddTerm((*test.RowBlock)(nil))
	//
	if types, e := types.GetTypes(); e != nil {
		t.Fatal(e)
	} else {
		m := Maker{types: types}
		desc, e := m.makeDesc("row_block", nil)
		require.NoError(t, e)
		expected := block.Dict{
			"message0": "%1",
			"type":     "row_block",
			"output":   nil,
			"args0": []block.Dict{
				{
					"name": "INPUT",
					"type": "input_value",
				},
			},
		}
		require.Equal(t, expected, desc)
	}
}
