package gblocks

import (
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/test"
	"github.com/stretchr/testify/require"
)

// TestConstraints - dependency pool generates constaints
func TestCheckNext(t *testing.T) {
	var types TypeCollector
	types.AddStatement((*test.CheckNext)(nil))
	if types, e := types.GetTypes(); e != nil {
		t.Fatal(e)
	} else {
		// when all types match, we should get an unlimited next statement
		// MOD: we can'd do that optimization b/c of mui inputs
		// the mui container gets an arbitrary set of mutating inputs
		// the blocks shouldnt swap between inputs, even though each isolated input looks like its unlimited
		m := Maker{types: types}
		desc, e := m.makeDesc("check_next", nil)
		require.NoError(t, e)
		expected := block.Dict{
			"message0":          "check next",
			"type":              "check_next",
			"tooltip":           "check next",
			"previousStatement": "check_next",
			//"nextStatement":     nil,
			"nextStatement": "check_next",
		}
		require.Equal(t, expected, desc)
	}
	types.AddStatement((*test.StackBlock)(nil))
	if types, e := types.GetTypes(); e != nil {
		t.Fatal(e)
	} else {
		// otherwise, a type limited next statement
		m := Maker{types: types}
		desc, e := m.makeDesc("check_next", nil)
		require.NoError(t, e)
		expected := block.Dict{
			"message0":          "check next",
			"type":              "check_next",
			"tooltip":           "check next",
			"previousStatement": "check_next",
			"nextStatement":     "check_next",
		}
		require.Equal(t, expected, desc)
	}
}

// TestConstraintsAny - dependency pool generates constaints
func TestConstraintsAny(t *testing.T) {
	var types TypeCollector
	types.AddTopStatement((*test.StackBlock)(nil))
	if types, e := types.GetTypes(); e != nil {
		t.Fatal(e)
	} else {
		m := Maker{types: types}
		desc, e := m.makeDesc("stack_block", nil)
		require.NoError(t, e)
		expected := block.Dict{
			"message0":      "stack block",
			"type":          "stack_block",
			"tooltip":       "stack block",
			"nextStatement": nil,
			// no prev because its a top block
		}
		require.Equal(t, expected, desc)
	}
}
