package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShapeDomSave(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e)
		//
		d := ws.GetDataById(b.Id).(*ShapeTest)
		d.Mutant = append(d.Mutant, &MutationEl{}, &MutationAlt{}, &MutationAlt{})
		//
		el := b.mutationToDom(ws)
		text := `<mutation><input name="MUTANT" types="mutation_el,mutation_alt" elements="0,1,1"></mutation>`
		require.Equal(t, text, el.OuterHTML())
	})
}

func TestShapeDomRestore(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e)
		//
		d := ws.GetDataById(b.Id).(*ShapeTest)
		d.Mutant = append(d.Mutant, &MutationEl{}, &MutationAlt{}, &MutationAlt{})
		//
		el := b.mutationToDom(ws)
		text := `<mutation><input name="MUTANT" types="mutation_el,mutation_alt" elements="0,1,1"></mutation>`
		require.Equal(t, text, el.OuterHTML())
	})
}
