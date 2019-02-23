package tools

import (
	. "github.com/ionous/gblocks/gtest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToolText(t *testing.T) {
	expected :=
		`<xml id="toolbox" style="display: none">` +
			/**/ `<category name="Logic" colour="%{BKY_LOGIC_HUE}">` +
			/* */ `<category name="If">` +
			/* */ `</category>` +
			/**/ `</category>` +
			`</xml>`
	elm := &Toolbox{
		Id:    "toolbox",
		Style: "display: none",
		Categories: Categories{
			&Category{
				Name: "Logic", Colour: "%{BKY_LOGIC_HUE}",
				Categories: Categories{
					&Category{Name: "If"},
				},
			},
		},
	}
	html := elm.OuterHTML()
	require.Equal(t, expected, html)
}

func TestToolStack(t *testing.T) {
	// text generated from blockly developer tools
	// https://blockly-demo.appspot.com/static/demos/blockfactory/index.html
	expected :=
		`<xml>` +
			`<block type="stack_block">` +
			/**/ `<next>` +
			/* */ `<block type="stack_block">` +
			/*   */ `<next>` +
			/*    */ `<block type="stack_block">` +
			/*    */ `</block>` +
			/*   */ `</next>` +
			/*  */ `</block>` +
			/* */ `</next>` +
			/**/ `</block>` +
			`</xml>`
	elm := NewToolbox(
		&StackBlock{
			NextStatement: &StackBlock{
				NextStatement: &StackBlock{
					NextStatement: nil,
				},
			},
		})
	require.Equal(t, expected, elm.OuterHTML())
}

func TestToolRow(t *testing.T) {
	expected :=
		`<xml>` +
			/* */ `<block type="row_block">` +
			/*  */ `<value name="INPUT">` +
			/*   */ `<block type="row_block">` +
			/*    */ `<value name="INPUT">` +
			/*     */ `<block type="row_block">` +
			/*     */ `</block>` +
			/*    */ `</value>` +
			/*   */ `</block>` +
			/*  */ `</value>` +
			/* */ `</block>` +
			`</xml>`
	elm := NewToolbox(
		&RowBlock{
			Input: &RowBlock{
				Input: &RowBlock{
					Input: nil,
				},
			},
		})
	require.Equal(t, expected, elm.OuterHTML())
}

func TestToolFieldBlock(t *testing.T) {
	expected :=
		`<xml>` +
			/**/ `<block type="field_block">` +
			/* */ `<field name="NUMBER">0</field>` +
			/**/ `</block>` +
			/**/ `<block type="field_block">` +
			/* */ `<field name="NUMBER">10</field>` +
			/**/ `</block>` +
			`</xml>`
	elm := NewToolbox(
		&FieldBlock{0}, &FieldBlock{10})
	require.Equal(t, expected, elm.OuterHTML())
}

func TestToolMutation(t *testing.T) {
	expected :=
		`<xml>` +
			/**/ `<block type="mutable_block">` +
			/* */ `<mutation>` +
			/*  */ `<atoms name="MUTANT">` +
			/*   */ `<atom type="atom_test"></atom>` +
			/*   */ `<atom type="atom_alt_test"></atom>` +
			/*   */ `<atom type="atom_test"></atom>` +
			/*  */ `</atoms>` +
			/* */ `</mutation>` + // the mutation itself has no field; the first dynamic element is 1, but is empty. so we start with MUTANT/2
			/* */ `<field name="MUTANT/2/ATOM_FIELD">Text</field>` +
			/* */ `<value name="MUTANT/3/ATOM_INPUT">` +
			/*  */ `<block type="mutable_block">` +
			/*   */ `<mutation>` +
			/*   */ `</mutation>` +
			/*  */ `</block>` +
			/* */ `</value>` +
			/**/ `</block>` +
			`</xml>`
	elm := NewToolbox(
		&MutableBlock{
			Mutant: TestMutation{&AtomTest{
				AtomInput: nil,
				NextStatement: &AtomAltTest{
					AtomField: "Text",
					NextStatement: &AtomTest{
						AtomInput:     &MutableBlock{},
						NextStatement: nil,
					}},
			}},
		})
	require.Equal(t, expected, elm.OuterHTML())
}

type StatementBlock struct {
	Do []interface{}
}

func TestToolStatement(t *testing.T) {
	expected :=
		`<xml>` +
			/**/ `<block type="statement_block">` +
			/* */ `<statement name="DO">` +
			/*   */ `<block type="stack_block">` +
			/*    */ `<next>` +
			/*     */ `<block type="stack_block">` +
			/*     */ `</block>` +
			/*   */ `</next>` +
			/*  */ `</block>` +
			/* */ `</statement>` +
			/**/ `</block>` +
			`</xml>`
		//
	elm := NewToolbox(
		&StatementBlock{
			Do: []interface{}{
				&StackBlock{},
				&StackBlock{},
			}})
	// t.Log(elm.OuterHTML())
	require.Equal(t, expected, elm.OuterHTML())
}

func TestToolEnum(t *testing.T) {
	expected :=
		`<xml>` +
			/**/ `<block type="enum_statement">` +
			/* */ `<field name="ENUM">AlternativeChoice</field>` +
			/**/ `</block>` +
			`</xml>`
		//
	elm := NewToolbox(
		&EnumStatement{Enum: AlternativeChoice})
	// t.Log(elm.OuterHTML())
	require.Equal(t, expected, elm.OuterHTML())
}
