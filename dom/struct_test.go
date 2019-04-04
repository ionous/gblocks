package dom

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

var common = struct {
	data Toolbox
	xml  string
}{
	xml: `<xml id="toolbox">` +
		/* */ `<category name="logic" colour="%{BKY_LOGIC_HUE}">` +
		/*  */ `<category name="If">` +
		/*   */ `<block type="row_block">` +
		/*   */ `</block>` +
		/*  */ `</category>` +
		/* */ `</category>` +
		/* */ `<block type="stack_block">` +
		/*  */ `<next>` +
		/*   */ `<block type="stack_block"></block>` +
		/*  */ `</next>` +
		/* */ `</block>` +
		/* */ `<block type="row_block">` +
		/*  */ `<value name="INPUT">` +
		/*   */ `<shadow type="row_block"></shadow>` +
		/*  */ `</value>` +
		/*  */ `<field name="NUMBER">10</field>` +
		/* */ `</block>` +
		/* */ `<block type="mutable_block">` +
		/*  */ `<mutation>` +
		/*   */ `<input name="MUTANT">` +
		/*    */ `<atom type="atom1"></atom>` +
		/*    */ `<atom type="atom2"></atom>` +
		/*   */ `</input>` +
		/*  */ `</mutation>` +
		/* */ `</block>` +
		`</xml>`,
	data: Toolbox{
		Id: "toolbox",
		Categories: Categories{
			&Category{
				Name:   "logic",
				Colour: "%{BKY_LOGIC_HUE}",
				Categories: Categories{
					&Category{
						Name: "If",
						Blocks: Blocks{
							&Block{Type: "row_block"},
						},
					},
				},
			},
		},
		Blocks: Blocks{
			&Block{
				Type: "stack_block",
				Next: ShapeLink{
					&Block{Type: "stack_block"},
				},
			},
			&Block{
				Type: "row_block",
				Items: ItemList{Items{
					&Value{
						Name: "INPUT",
						Input: ShapeInput{
							&Shadow{
								Type: "row_block",
							}},
					},
					&Field{
						Name:    "NUMBER",
						Content: "10",
					},
				}},
			},
			&Block{
				Type: "mutable_block",
				Mutations: &Mutations{
					&Mutation{
						Input: "MUTANT",
						Atoms: Atoms{
							[]string{"atom1", "atom2"},
						},
					},
				},
			},
		},
	},
}

func TestMarshal(t *testing.T) {
	output, e := common.data.MarshalToString()
	require.NoError(t, e)
	require.Equal(t, common.xml, string(output))
}

func TestUnmarshalIndent(t *testing.T) {
	data, e := NewToolboxFromString(common.xml)
	require.NoError(t, e)
	// t.Log(pretty.Sprint(data))
	if diff := pretty.Diff(common.data, *data); len(diff) != 0 {
		t.Fatal(diff)
	}
}
