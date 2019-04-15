package tin

import (
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
	"github.com/ionous/gblocks/test"
	"github.com/stretchr/testify/require"
)

// verify the mui block names and descriptions generated for test.BlockMutation
func TestQuarks(t *testing.T) {
	var ms Mutables
	mut, e := ms.addMutation((*test.BlockMutation)(nil),
		(*test.AtomTest)(nil),
		(*test.AtomAltTest)(nil))
	require.NoError(t, e)
	//
	quarks := mutant.PaletteQuarks(mut)
	expectedQuarks := []string{"mui$block_mutation$atom_test", "mui$block_mutation$atom_alt_test"}
	require.Equal(t, expectedQuarks, quarks)
	var pal []interface{}
	for q, ok := mut.Quarks(false); ok; q, ok = q.NextQuark() {
		pal = append(pal, mutant.DescribeQuark(q))
	}

	expected := []interface{}{
		block.Dict{
			"type":              "mui$block_mutation$",
			"message0":          "block",
			"previousStatement": "block_mutation$",
			"nextStatement":     []string{"atom_test", "atom_alt_test"},
		},
		block.Dict{
			"type":              "mui$block_mutation$atom_test",
			"message0":          "atom test",
			"previousStatement": "atom_test",
			"nextStatement":     []string{"atom_test", "atom_alt_test"},
		},
		block.Dict{
			"type":              "mui$block_mutation$atom_alt_test",
			"message0":          "atom alt test",
			"previousStatement": "atom_alt_test",
			"nextStatement":     []string{"atom_test", "atom_alt_test"},
		},
	}
	require.Equal(t, expected, pal)

	mins := mutant.InMutations{
		[]string{"IN"},
		map[string]mutant.InMutation{"IN": mut},
	}
	container := mins.DescribeContainer(mutant.ContainerName("block"))
	expectedContainer := block.Dict{
		"type":     "mui$block",
		"message0": "%1",
		"args0": []block.Dict{
			{
				"name":  "IN",
				"type":  "input_statement",
				"check": "block_mutation$",
			},
		},
	}
	require.Equal(t, expectedContainer, container)
}

func TestQuarksNoFixedFields(t *testing.T) {
	type BlockMutation struct {
		NextStatement test.NextAtom
	}
	//
	var ms Mutables
	mut, e := ms.addMutation((*BlockMutation)(nil),
		(*test.AtomTest)(nil),
		(*test.AtomAltTest)(nil))
	require.NoError(t, e)
	//
	names := mutant.PaletteQuarks(mut)
	expectedQuarks := []string{"mui$block_mutation$atom_test", "mui$block_mutation$atom_alt_test"}
	require.Equal(t, expectedQuarks, names)
}

func TestQuarksEmpty(t *testing.T) {
	type BlockMutation struct{}
	//
	var ms Mutables
	mut, e := ms.addMutation((*BlockMutation)(nil))
	require.NoError(t, e)
	//
	names := mutant.PaletteQuarks(mut)
	require.Empty(t, names)
}
