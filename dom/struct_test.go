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
	xml: `<xml>` +
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
		/*   */ `<pin name="MUTANT">` +
		/*    */ `<atom name="name1" type="atom1"></atom>` +
		/*    */ `<atom name="name2" type="atom2"></atom>` +
		/*   */ `</pin>` +
		/*  */ `</mutation>` +
		/* */ `</block>` +
		`</xml>`,
	data: Toolbox{
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
				Mutation: &BlockMutation{Mutations{
					&Mutation{
						Input: "MUTANT",
						Atoms: Atoms{
							&Atom{"name1", "atom1"},
							&Atom{"name2", "atom2"},
						},
					},
				}},
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

// BlockMutation are directly parsed during mui compose/decompose
func TestUnmarshalMutation(t *testing.T) {
	str := `` +
		/*  */ `<mutation>` +
		/*   */ `<pin name="A">` +
		/*    */ `<atom name="one" type="atom"/>` +
		/*   */ `</pin>` +
		/*   */ `<pin name="B">` +
		/*    */ `<atom name="two" type="atom"/>` +
		/*   */ `</pin>` +
		/*  */ `</mutation>`

	if ms, e := UnmarshalMutation(str); e != nil {
		t.Fatal(e)
	} else {
		expected := BlockMutation{
			Mutations{
				&Mutation{
					Input: "A",
					Atoms: Atoms{
						&Atom{"one", "atom"},
					},
				},
				&Mutation{
					Input: "B",
					Atoms: Atoms{
						&Atom{"two", "atom"},
					},
				},
			},
		}
		require.Equal(t, expected, ms)
	}
}
