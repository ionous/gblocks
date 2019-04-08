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
						Blocks: BlockList{Blocks{
							&Block{Type: "row_block"},
						}},
					},
				},
			},
		},
		Blocks: Blocks{
			&Block{
				Type: "stack_block",
				Next: BlockLink{
					&Block{Type: "stack_block"},
				},
			},
			&Block{
				Type: "row_block",
				Items: ItemList{Items{
					&Value{
						Name: "INPUT",
						Input: BlockInput{
							&Block{
								Type:   "row_block",
								Shadow: true,
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
	var html string
	require.NotPanics(t, func() {
		html = common.data.OuterHTML()
	})
	require.Equal(t, common.xml, html)

}

func TestUnmarshalIndent(t *testing.T) {
	data, e := NewToolboxFromString(common.xml)
	require.NoError(t, e)
	// t.Log(pretty.Sprint(data))
	if diff := pretty.Diff(common.data, *data); len(diff) != 0 {
		t.Fatal(diff)
	}
}
