package gblocks

import (
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
	tool := NewTools(nil)
	elm := tool.Box("xml", Attrs{"id": "toolbox", "style": "display: none"},
		tool.Box("category", Attrs{"name": "Logic", "colour": "%{BKY_LOGIC_HUE}"},
			tool.Box("category", Attrs{"name": "If"}),
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
			/*    */ `<block type="stack_block">` +
			/*    */ `</block>` +
			/*   */ `</next>` +
			/*  */ `</block>` +
			/* */ `</next>` +
			/**/ `</block>` +
			`</xml>`
	elm := NewToolData(NewDomElement("xml"),
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
	elm := NewToolData(NewDomElement("xml"),
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
	elm := NewToolData(NewDomElement("xml"),
		&FieldBlock{0}, &FieldBlock{10})
	require.Equal(t, expected, elm.OuterHTML())
}

func TestToolMutation(t *testing.T) {
	expected :=
		`<xml>` +
			/**/ `<block type="shape_test">` +
			/* */ `<mutation>` +
			/*  */ `<input name="MUTANT" types="mutation_el,mutation_alt" elements="0,1">` +
			/* */ `</mutation>` +
			/**/ `</block>` +
			`</xml>`
	elm := NewToolData(
		NewDomElement("xml"),
		&ShapeTest{Mutant: []interface{}{
			&MutationEl{},
			&MutationAlt{},
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
			/*     */ `<block type="stack_block">` +
			/*     */ `</block>` +
			/*   */ `</next>` +
			/*  */ `</block>` +
			/* */ `</statement>` +
			/**/ `</block>` +
			`</xml>`
		//
	elm := NewToolData(NewDomElement("xml"),
		&StatementBlock{
			Do: []interface{}{
				&StackBlock{},
				&StackBlock{},
			}})
	// t.Log(elm.OuterHTML())
	require.Equal(t, expected, elm.OuterHTML())
}
