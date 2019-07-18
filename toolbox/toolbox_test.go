package toolbox_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/ionous/gblocks/block"
	. "github.com/ionous/gblocks/test"
	"github.com/ionous/gblocks/toolbox"
	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	// text generated from blockly developer tools
	// https://blockly-demo.appspot.com/static/demos/blockfactory/index.html
	expected := `` +
		`<block type="stack_block">` +
		/**/ `<next>` +
		/* */ `<block type="stack_block">` +
		/*  */ `<next>` +
		/*   */ `<block type="stack_block">` +
		/*   */ `</block>` +
		/*  */ `</next>` +
		/* */ `</block>` +
		/**/ `</next>` +
		`</block>`

	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, types).
		AddStatement(
			&StackBlock{
				NextStatement: &StackBlock{
					NextStatement: &StackBlock{
						NextStatement: nil,
					},
				},
			}).
		Blocks()
	if collected, e := types.Collected(); e != nil {
		t.Fatal(e)
	} else {
		expected := []string{"stack_block (MidBlock)"}
		require.Equal(t, expected, collected)
	}
	html, e := xml.Marshal(blocks)
	require.NoError(t, e)
	require.Equal(t, expected, string(html))
}

func TestRow(t *testing.T) {
	expected := `` +
		/* */ `<block type="row_block">` +
		/*  */ `<value name="INPUT">` +
		/*   */ `<block type="row_block">` +
		/*    */ `<value name="INPUT">` +
		/*     */ `<block type="row_block">` +
		/*     */ `</block>` +
		/*    */ `</value>` +
		/*   */ `</block>` +
		/*  */ `</value>` +
		/* */ `</block>`

	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, types).
		AddTerm(
			&RowBlock{
				Input: &RowBlock{
					Input: &RowBlock{
						Input: nil,
					},
				},
			}).
		Blocks()
	if collected, e := types.Collected(); e != nil {
		t.Fatal(e)
	} else {
		expected := []string{"row_block (TermBlock)"}
		require.Equal(t, expected, collected)
	}
	html, e := xml.Marshal(blocks)
	require.NoError(t, e)
	require.Equal(t, expected, string(html))
}

// test generation of blocks containing fields
// test no ids and no collector while we are at it.
func TestFieldBlock(t *testing.T) {
	expected := `` +
		/**/ `<block type="field_block">` +
		/**/ `</block>` +
		/**/ `<block type="field_block">` +
		/* */ `<field name="NUMBER">10</field>` +
		/**/ `</block>`

	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, types).
		AddTerm(&FieldBlock{0}).
		AddTerm(&FieldBlock{10}).
		Blocks()
	if collected, e := types.Collected(); e != nil {
		t.Fatal(e)
	} else {
		expected := []string{"field_block (TermBlock)"}
		require.Equal(t, expected, collected)
	}

	html, e := xml.Marshal(blocks)
	require.NoError(t, e)
	require.Equal(t, expected, string(html))
}

// test no ids and no collector while we are at it.
func TestNoCollection(t *testing.T) {
	expected := `` +
		/**/ `<block type="field_block">` +
		/**/ `</block>` +
		/**/ `<block type="field_block">` +
		/* */ `<field name="NUMBER">10</field>` +
		/**/ `</block>`
	blocks := toolbox.NewBlocks(toolbox.NoShadow, nil).
		AddTerm(&FieldBlock{0}).
		AddTerm(&FieldBlock{10}).
		Blocks()
	html, e := xml.Marshal(blocks)
	require.NoError(t, e)
	require.Equal(t, expected, string(html))
}

func TestMutations(t *testing.T) {
	type MutationExtraless struct {
		NextStatement NextAtom
	}
	type BlockExtraless struct {
		Mutant *MutationExtraless `input:"mutation"`
	}

	// "a, INPUT_NAME, atomNum, FIELD_NAME"
	el := func(tag string, parts ...string) string {
		return fmt.Sprintf("<%s name=%q>", tag, block.Scope(parts...))
	}
	expected := `` +
		/**/ `<block type="mutable_block">` +
		/* */ `<mutation>` +
		/*  */ `<pin name="MUTANT">` +
		/*   */ `<atom type="block_mutation"></atom>` +
		/*   */ `<atom type="atom_test"></atom>` +
		/*   */ `<atom type="atom_alt_test"></atom>` +
		/*   */ `<atom type="atom_test"></atom>` +
		/*  */ `</pin>` +
		/* */ `</mutation>` +
		/* */ el("field", "a", "MUTANT", "0", "ATOM_FIELD") + `Text</field>` +
		/* */ el("value", "a", "MUTANT", "1", "ATOM_INPUT") +
		/*  */ `<block type="mutable_block">` +
		/*   */ `<mutation>` +
		/*    */ `<pin name="MUTANT">` +
		/*     */ `<atom type="block_mutation"></atom>` +
		/*  */ `</pin>` +
		/*   */ `</mutation>` +
		/*  */ `</block>` +
		/* */ `</value>` +
		/**/ `</block>` +
		/**/ `<block type="block_extraless">` +
		/* */ `<mutation>` +
		/*  */ `<pin name="MUTANT">` +
		/*   */ `<atom type="atom_test"></atom>` +
		/*  */ `</pin>` +
		/* */ `</mutation>` +
		/**/ `</block>` +
		/**/ `<block type="mutable_block">` +
		/* */ `<mutation>` +
		/*  */ `<pin name="MUTANT">` +
		/*   */ `<atom type="block_mutation"></atom>` +
		/*  */ `</pin>` +
		/* */ `</mutation>` +
		/**/ `</block>` +
		/**/ `<block type="block_extraless">` +
		/* */ `<mutation>` +
		/* */ `</mutation>` +
		/**/ `</block>`

	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, types).
		AddTerm(
			&MutableBlock{
				Mutant: &BlockMutation{
					NextStatement: &AtomTest{ // atom 0
						AtomInput: nil,
						NextStatement: &AtomAltTest{ // atom 1
							AtomField: "Text",
							NextStatement: &AtomTest{ // atom 2
								AtomInput:     &MutableBlock{},
								NextStatement: nil,
							}},
					}},
			}).
		// one atom: the atom test
		AddTerm(
			&BlockExtraless{
				Mutant: &MutationExtraless{
					NextStatement: &AtomTest{ // atom 0
					}},
			}).
		// one atom: the mutation itself
		AddTerm(&MutableBlock{}).
		// no atoms at all
		AddTerm(&BlockExtraless{}).
		Blocks()
	if collected, e := types.Collected(); e != nil {
		t.Fatal(e)
	} else {
		// note: the atoms arent registered; they need to be explicitly registered
		// and that's not part of toolbox; its a separate sort of toolbox.
		expected := []string{"mutable_block (TermBlock)", "block_extraless (TermBlock)"}
		require.Equal(t, expected, collected)
	}
	html, e := xml.Marshal(blocks)
	require.NoError(t, e)
	require.Equal(t, expected, string(html))
}

type StatementBlock struct {
	Do interface{} `input:"statement"`
}

func TestStatement(t *testing.T) {
	expected := `` +
		/**/ `<block type="statement_block">` +
		/* */ `<statement name="DO">` +
		/*   */ `<shadow type="stack_block">` +
		/*    */ `<next>` +
		/*     */ `<shadow type="stack_block">` +
		/*     */ `</shadow>` +
		/*   */ `</next>` +
		/*  */ `</shadow>` +
		/* */ `</statement>` +
		/**/ `</block>`

	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.SubShadow, types).
		AddStatement(
			&StatementBlock{
				Do: &StackBlock{NextStatement: &StackBlock{}},
			}).
		Blocks()
	if collected, e := types.Collected(); e != nil {
		t.Fatal(e)
	} else {
		expected := []string{"statement_block (MidBlock)", "stack_block (MidBlock)"}
		require.Equal(t, expected, collected)
	}
	html, e := xml.Marshal(blocks)
	require.NoError(t, e)
	require.Equal(t, expected, string(html))
}

func TestEnum(t *testing.T) {
	expected := `` +
		/**/ `<block type="enum_statement">` +
		/* */ `<field name="ENUM">AlternativeChoice</field>` +
		/**/ `</block>`
	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, types).
		AddStatement(&EnumStatement{AlternativeChoice}).
		Blocks()
	if collected, e := types.Collected(); e != nil {
		t.Fatal(e)
	} else {
		expected := []string{"enum_statement (MidBlock)"}
		require.Equal(t, expected, collected)
	}
	html, e := xml.Marshal(blocks)
	require.NoError(t, e)
	require.Equal(t, expected, string(html))
}
