package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToolText(t *testing.T) {
	expected :=
		`<xml id="toolbox" style="display: none">` +
			/**/ `<category name="Logic" colour="%{BKY_LOGIC_HUE}">` +
			/* */ `<category name="If"/>` +
			/**/ `</category>` +
			`</xml>`
	elm := Toolbox("xml", Attrs{"id": "toolbox", "style": "display: none"},
		Toolbox("category", Attrs{"name": "Logic", "colour": "%{BKY_LOGIC_HUE}"},
			Toolbox("category", Attrs{"name": "If"}),
		),
	)
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
			/*    */ `<block type="stack_block"/>` +
			/*   */ `</next>` +
			/*  */ `</block>` +
			/* */ `</next>` +
			/**/ `</block>` +
			`</xml>`
	elm := Toolbox("xml", nil,
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
			/*     */ `<block type="row_block"/>` +
			/*    */ `</value>` +
			/*   */ `</block>` +
			/*  */ `</value>` +
			/* */ `</block>` +
			`</xml>`
	elm := Toolbox("xml", nil,
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
	elm := Toolbox("xml", nil,
		&FieldBlock{0}, &FieldBlock{10})
	require.Equal(t, expected, elm.OuterHTML())
}

func TestToolMutation(t *testing.T) {
	expected :=
		`<xml>` +
			/**/ `<block type="shape_test">` +
			/* */ `<mutation>` +
			/*  */ `<atoms name="MUTANT">` +
			/*   */ `<atom type="atom_test"/>` +
			/*   */ `<atom type="atom_alt_test"/>` +
			/*  */ `</atoms>` +
			/* */ `</mutation>` +
			/* */ `<statement name="MUTANT">` +
			/*  */ `<block type="atom_test">` +
			/*   */ `<next>` +
			/*    */ `<block type="atom_alt_test"/>` +
			/*   */ `</next>` +
			/*  */ `</block>` +
			/* */ `</statement>` +
			/**/ `</block>` +
			`</xml>`
	elm := NewTools(
		NewXmlElement("xml"),
		&ShapeTest{Mutant: []interface{}{
			&AtomTest{},
			&AtomAltTest{},
		}})
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
			/*     */ `<block type="stack_block"/>` +
			/*   */ `</next>` +
			/*  */ `</block>` +
			/* */ `</statement>` +
			/**/ `</block>` +
			`</xml>`
		//
	elm := Toolbox("xml", nil,
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
	elm := Toolbox("xml", nil,
		&EnumStatement{Enum: AlternativeChoice})
	// t.Log(elm.OuterHTML())
	require.Equal(t, expected, elm.OuterHTML())
}
