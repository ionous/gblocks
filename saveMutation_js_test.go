package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShapeDomSave(t *testing.T) {
	testShape(t, func(ws *Workspace, reg *Registry) {
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e, "created block")
		require.Equal(t, 3, b.NumInputs(), "initial inputs")
		//
		if in, index := b.InputByName("MUTANT"); index < 0 {
			t.Fatal("missing input")
		} else if m := in.Mutation(); m == nil {
			t.Fatal("missing mutation")
		} else {
			for i, atomType := range []TypeName{"atom_test", "atom_alt_test", "atom_alt_test"} {
				numInputs, e := m.addAtom(reg, atomType)
				require.NoError(t, e, "added atom", i)
				require.Equal(t, 1, numInputs, "added inputs", i)
			}
		}
		require.Equal(t, b.NumInputs(), 6, "expanded inputs")
		//
		el := b.mutationToDom()
		text := `<mutation><atoms name="MUTANT">` +
			/**/ `<atom type="atom_test"></atom>` +
			/**/ `<atom type="atom_alt_test"></atom>` +
			/**/ `<atom type="atom_alt_test"></atom>` +
			`</atoms></mutation>`
		require.Equal(t, text, el.OuterHTML())
	})
}

// func xTestShapeDomRestore(t *testing.T) {
// 	testShape(t, func(ws *Workspace, reg*Registry) {
// 		b, e := ws.NewBlock((*ShapeTest)(nil))
// 		require.NoError(t, e)
// 		//
// 		d := ws.GetDataById(b.Id).(*ShapeTest)
// 		d.Mutant = append(d.Mutant, &AtomTest{}, &AtomAltTest{}, &AtomAltTest{})
// 		//
// 		el := b.mutationToDom(ws)
// 		text := `<mutation><data name="MUTANT" types="mutation_el,mutation_alt" elements="0,1,1"></data></mutation>`
// 		require.Equal(t, text, el.OuterHTML())
// 	})
// }
