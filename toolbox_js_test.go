package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func xTestToolText(t *testing.T) {
	expected :=
		`<xml id="toolbox" style="display: none">` +
			/**/ `<category name="Logic" colour="%{BKY_LOGIC_HUE}">` +
			/* */ `<category name="If">` +
			/* */ `</category>` +
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

func xTestToolStack(t *testing.T) {
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

func xTestToolRow(t *testing.T) {
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

func xTestToolFieldBlock(t *testing.T) {
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

func xTestToolMutation(t *testing.T) {
	expected :=
		`<xml>` +
			/**/ `<block type="shape_test">` +
			/* */ `<mutation>` +
			/*  */ `<data name="MUTANT" types="mutation_el,mutation_alt" elements="0,1">` +
			/*  */ `</data>` +
			/* */ `</mutation>` +
			/* */ `<statement name="MUTANT">` +
			/*  */ `<block type="mutation_el">` +
			/*   */ `<next>` +
			/*    */ `<block type="mutation_alt">` +
			/*    */ `</block>` +
			/*   */ `</next>` +
			/*  */ `</block>` +
			/* */ `</statement>` +
			/**/ `</block>` +
			`</xml>`
	elm := NewTools(
		NewDomElement("xml"),
		&ShapeTest{Mutant: []interface{}{
			&AtomTest{},
			&AtomAltTest{},
		}})
	require.Equal(t, expected, elm.OuterHTML())
}

type StatementBlock struct {
	Do []interface{}
}

func xTestToolStatement(t *testing.T) {
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
	elm := Toolbox("xml", nil,
		&StatementBlock{
			Do: []interface{}{
				&StackBlock{},
				&StackBlock{},
			}})
	// t.Log(elm.OuterHTML())
	require.Equal(t, expected, elm.OuterHTML())
}

func xTestToolEnum(t *testing.T) {
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
