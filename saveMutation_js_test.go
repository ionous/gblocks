package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func expectedDom() *XmlElement {
	return Toolbox("mutation", nil,
		Toolbox("atoms", Attrs{"name": "MUTANT"},
			Toolbox("atom", Attrs{"type": "atom_test"}),
			Toolbox("atom", Attrs{"type": "atom_alt_test"}),
			Toolbox("atom", Attrs{"type": "atom_test"}),
		))
}

func TestShapeDom(t *testing.T) {
	text := `` +
		`<mutation>` +
		/**/ `<atoms name="MUTANT">` +
		/* */ `<atom type="atom_test"/>` +
		/* */ `<atom type="atom_alt_test"/>` +
		/* */ `<atom type="atom_test"/>` +
		/**/ `</atoms>` +
		`</mutation>`
	require.Equal(t, text, expectedDom().OuterHTML())
}

func TestShapeAddAtom(t *testing.T) {
	testShape(t, func(ws *Workspace, reg *Registry) {
		b, e := ws.NewBlock((*ShapeTest)(nil))
		// note: unfortunately, fields dont keep their names; so the last input is nil
		// ( see Blockly.Block.prototype.interpolate_ )
		require.NoError(t, e, "created block")
		require.Equal(t, 3, b.NumInputs(), "initial inputs")
		expected := []string{"INPUT", "MUTANT", "<field>"}
		require.Equal(t, expected, reduceInputs(b))
		//
		if in, index := b.InputByName("MUTANT"); index < 0 {
			t.Fatal("missing input")
		} else if m := in.Mutation(); m == nil {
			t.Fatal("missing mutation")
		} else {
			t.Log("first")
			{
				numInputs, e := m.addAtom(reg, "atom_test")
				require.NoError(t, e, "added atom", 1)
				require.Equal(t, 1, numInputs, "added inputs")
				//
				require.Equal(t, b.NumInputs(), 4, "expanded inputs")
				expected := []string{"INPUT", "MUTANT", "MUTANT/0/SUB_INPUT", "<field>"}
				require.Equal(t, expected, reduceInputs(b))
			}
			t.Log("second")
			{
				numInputs, e := m.addAtom(reg, "atom_alt_test")
				require.NoError(t, e, "added atom", 1)
				require.Equal(t, 1, numInputs, "added inputs")
				//
				require.Equal(t, b.NumInputs(), 5, "expanded inputs")
				expected := []string{"INPUT", "MUTANT", "MUTANT/0/SUB_INPUT", "<field>", "<field>"}
				require.Equal(t, expected, reduceInputs(b))
			}
			t.Log("third")
			{
				numInputs, e := m.addAtom(reg, "atom_test")
				require.NoError(t, e, "added atom", 1)
				require.Equal(t, 1, numInputs, "added inputs")
				//
				// note: unfortunately, fields dont keep their names; so the last input is nil
				// ( see Blockly.Block.prototype.interpolate_ )
				require.Equal(t, b.NumInputs(), 6, "expanded inputs")
				expected := []string{"INPUT", "MUTANT", "MUTANT/0/SUB_INPUT", "<field>", "MUTANT/2/SUB_INPUT", "<field>"}
				require.Equal(t, expected, reduceInputs(b))
			}
		}
	})
}

func TestShapeSave(t *testing.T) {
	testShape(t, func(ws *Workspace, reg *Registry) {
		b, e := buildMutation(ws, reg, t)
		require.NoError(t, e)
		//
		el := b.mutationToDom()
		require.Equal(t, expectedDom().OuterHTML(), el.OuterHTML())
	})
}

func TestShapeRestore(t *testing.T) {
	testShape(t, func(ws *Workspace, reg *Registry) {
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e)
		//
		added, e := b.domToMutation(reg, expectedDom())
		require.NoError(t, e)
		require.Equal(t, 3, added)
		//
		if in, index := b.InputByName("MUTANT"); index < 0 {
			t.Fatal("missing input")
		} else if m := in.Mutation(); m == nil {
			t.Fatal("missing mutation")
		} else {
			expected := []string{"INPUT", "MUTANT", "MUTANT/0/SUB_INPUT", "<field>", "MUTANT/2/SUB_INPUT", "<field>"}
			require.Equal(t, expected, reduceInputs(b))
		}
	})
}
