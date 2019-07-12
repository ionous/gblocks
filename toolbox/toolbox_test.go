package toolbox_test

import (
	"encoding/xml"
	"testing"

	. "github.com/ionous/gblocks/test"
	"github.com/ionous/gblocks/toolbox"
	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	// text generated from blockly developer tools
	// https://blockly-demo.appspot.com/static/demos/blockfactory/index.html
	expected := `` +
		`<block id="bl1" type="stack_block">` +
		/**/ `<next>` +
		/* */ `<block id="bl2" type="stack_block">` +
		/*  */ `<next>` +
		/*   */ `<block id="bl3" type="stack_block">` +
		/*   */ `</block>` +
		/*  */ `</next>` +
		/* */ `</block>` +
		/**/ `</next>` +
		`</block>`
	ids := &toolbox.IdGenerator{}
	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, ids, types).
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
		/* */ `<block id="bl1" type="row_block">` +
		/*  */ `<value name="INPUT">` +
		/*   */ `<block id="bl2" type="row_block">` +
		/*    */ `<value name="INPUT">` +
		/*     */ `<block id="bl3" type="row_block">` +
		/*     */ `</block>` +
		/*    */ `</value>` +
		/*   */ `</block>` +
		/*  */ `</value>` +
		/* */ `</block>`
	ids := &toolbox.IdGenerator{}
	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, ids, types).
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
		/**/ `<block id="bl1" type="field_block">` +
		/**/ `</block>` +
		/**/ `<block id="bl2" type="field_block">` +
		/* */ `<field name="NUMBER">10</field>` +
		/**/ `</block>`
	ids := &toolbox.IdGenerator{}
	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, ids, types).
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
	blocks := toolbox.NewBlocks(toolbox.NoShadow, nil, nil).
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

	// "a$ wsBlockId $ INPUT_NAME $ atomNum $ FIELD_NAME"
	expected := `` +
		/**/ `<block id="bl1" type="mutable_block">` +
		/* */ `<mutation>` +
		/*  */ `<pin name="MUTANT">` +
		/*   */ `<atom type="block_mutation$"></atom>` +
		/*   */ `<atom type="atom_test"></atom>` +
		/*   */ `<atom type="atom_alt_test"></atom>` +
		/*   */ `<atom type="atom_test"></atom>` +
		/*  */ `</pin>` +
		/* */ `</mutation>` +
		/* */ `<field name="a$bl1$MUTANT$0$ATOM_FIELD">Text</field>` +
		/* */ `<value name="a$bl1$MUTANT$1$ATOM_INPUT">` +
		/*  */ `<block id="bl2" type="mutable_block">` +
		/*   */ `<mutation>` +
		/*    */ `<pin name="MUTANT">` +
		/*     */ `<atom type="block_mutation$"></atom>` +
		/*  */ `</pin>` +
		/*   */ `</mutation>` +
		/*  */ `</block>` +
		/* */ `</value>` +
		/**/ `</block>` +
		/**/ `<block id="bl3" type="block_extraless">` +
		/* */ `<mutation>` +
		/*  */ `<pin name="MUTANT">` +
		/*   */ `<atom type="atom_test"></atom>` +
		/*  */ `</pin>` +
		/* */ `</mutation>` +
		/**/ `</block>` +
		/**/ `<block id="bl4" type="mutable_block">` +
		/* */ `<mutation>` +
		/*  */ `<pin name="MUTANT">` +
		/*   */ `<atom type="block_mutation$"></atom>` +
		/*  */ `</pin>` +
		/* */ `</mutation>` +
		/**/ `</block>` +
		/**/ `<block id="bl5" type="block_extraless">` +
		/* */ `<mutation>` +
		/* */ `</mutation>` +
		/**/ `</block>`

	ids := &toolbox.IdGenerator{}
	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.NoShadow, ids, types).
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
		/**/ `<block id="bl1" type="statement_block">` +
		/* */ `<statement name="DO">` +
		/*   */ `<shadow id="bl2" type="stack_block">` +
		/*    */ `<next>` +
		/*     */ `<shadow id="bl3" type="stack_block">` +
		/*     */ `</shadow>` +
		/*   */ `</next>` +
		/*  */ `</shadow>` +
		/* */ `</statement>` +
		/**/ `</block>`
	ids := &toolbox.IdGenerator{}
	types := &testCollector{}
	blocks := toolbox.NewBlocks(toolbox.SubShadow, ids, types).
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
	blocks := toolbox.NewBlocks(toolbox.NoShadow, nil, types).
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
