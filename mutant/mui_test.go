package mutant_test

import (
	"sort"
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mock"
	"github.com/ionous/gblocks/mutant"
	"github.com/stretchr/testify/require"
)

var common = struct {
	inputs, quarks []string
	inputAtoms     mutant.MutableInputs
	muiContainer   block.Dict
	atomProducts   map[string][]mock.MockAtom
	// workspace inputs after expanding the mutable input atoms
	expandedInputs []string
}{
	inputs: []string{"I1:input_dummy", "I2:input_dummy", "I3:input_dummy"},
	quarks: []string{"a1", "a2", "a3"},
	inputAtoms: mutant.MutableInputs{
		"I1": []string{"a1", "a1"},
		"I2": []string{"a2"},
		"I3": []string{"a1", "a2", "a3"},
	},
	muiContainer: block.Dict{
		"type":     "mui$test",
		"message0": "%1 %2 %3",
		"args0": []block.Dict{{
			"name":  "I1",
			"type":  "input_statement",
			"check": "a1",
		}, {
			"type":  "input_statement",
			"name":  "I2",
			"check": "a1",
		}, {
			"name":  "I3",
			"type":  "input_statement",
			"check": "a1",
		}},
	},
	atomProducts: map[string][]mock.MockAtom{
		"a1": []mock.MockAtom{{Name: "TERM", Type: block.ValueInput}},
		"a2": []mock.MockAtom{{Name: "NUM", Type: block.NumberField}},
		"a3": []mock.MockAtom{{Name: "TEXT", Type: block.TextField}, {Name: "STATE", Type: block.StatementInput}},
	},
	expandedInputs: []string{
		// a$ blockId $ INPUT $ atomNum $ FIELD : type
		"I1:input_dummy",
		/*a1*/ "a$mock$I1$0$TERM:input_value",
		/*a1*/ "a$mock$I1$1$TERM:input_value",
		"I2:input_dummy",
		/*a2*/ "a$mock$I2$0$NUM:field_number",
		"I3:input_dummy",
		/*a1*/ "a$mock$I3$0$TERM:input_value",
		/*a2*/ "a$mock$I3$1$NUM:field_number",
		/*a3*/ "a$mock$I3$2$TEXT:field_input", "a$mock$I3$2$STATE:input_statement",
	},
}

func TestCreateMockBlock(t *testing.T) {
	// create a workspace block where we will be writing data to
	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
	require.Equal(t, b.NumInputs(), 3)
	require.Equal(t, common.inputs, listInputs(b))
}

func TestDescribeContainer(t *testing.T) {
	mutator := mock.NewInMutations(common.inputs, common.quarks)
	muiContainer := mutator.DescribeContainer(mutant.ContainerName("test"))
	require.Equal(t, common.muiContainer, muiContainer)
}

func TestPreregister(t *testing.T) {
	var reg mock.Registry
	mutator := mock.NewInMutations(common.inputs, common.quarks)
	require.NoError(t, mutator.Preregister("test", &reg))
	var keys []string
	for k, _ := range reg.Blocks {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// we expect to see unique mutations for each input
	// ( named after its input )
	// and we expect to see unique quarks for each mutation
	// ( also the mui container which can hold the quarks )
	expected := []string{
		"mui$I1$a1",
		"mui$I1$a2",
		"mui$I1$a3",
		"mui$I2$a1",
		"mui$I2$a2",
		"mui$I2$a3",
		"mui$I3$a1",
		"mui$I3$a2",
		"mui$I3$a3",
		"mui$test"}
	require.Equal(t, expected, keys)
}

func TestCreateMui(t *testing.T) {
	var reg mock.Registry
	muispace := reg.NewMockSpace()
	mutator := mock.NewInMutations(common.inputs, common.quarks)
	require.NoError(t, mutator.Preregister("mockType", &reg))
	wsblock := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
	c, e := mutator.CreateMui(muispace, wsblock, common.inputAtoms)
	require.NoError(t, e, "CreateMui")
	idTypes := listStack(c)
	expectedIdTypes := []string{
		"mock$I1$0:mui$I1$a1", "mock$I1$1:mui$I1$a1",
		"mock$I2$0:mui$I2$a2",
		"mock$I3$0:mui$I3$a1", "mock$I3$1:mui$I3$a2", "mock$I3$2:mui$I3$a3",
	}
	require.Equal(t, expectedIdTypes, idTypes, "expectedIdTypes")
}

func listStack(b block.Shape) (ret []string) {
	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
		in := b.Input(i)
		block.VisitStack(in, func(b block.Shape) bool {
			ret = append(ret, b.(interface{ String() string }).String())
			return true // keepGoing
		})
	}
	return
}

// note: with the mock implementation, every field gets promoted to a connectionless input.
func listInputs(b block.Shape) (ret []string) {
	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
		in := b.Input(i)
		ret = append(ret, in.(interface{ String() string }).String())
	}
	return
}

func TestDistillMui(t *testing.T) {
	var reg mock.Registry
	// create a workspace block where we will be writing data to
	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
	// first create a container and give it some fake atom data.
	muispace := reg.NewMockSpace()
	// have to create the mui in order to fill from it.
	// (  a container with inputs containing stacks of atoms )
	mutator := mock.NewInMutations(common.inputs, common.quarks)
	require.NoError(t, mutator.Preregister("mockType", &reg))
	c, e := mutator.CreateMui(muispace, b, common.inputAtoms)
	require.NoError(t, e)
	// now, fill the b with the blocks from the container
	db := &mock.MockDatabase{common.atomProducts}
	inputAtoms, e := mutator.DistillMui(b, c, db, nil)
	require.NoError(t, e)
	require.Equal(t, common.inputAtoms, inputAtoms)
	//
	expandedInputs := listInputs(b)
	require.Equal(t, common.expandedInputs, expandedInputs, "expandedInputs")
}

// start a block with atoms then remove those atoms and ensure the block is "empty"
func TestRemoveAtoms(t *testing.T) {
	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.expandedInputs))
	expandedInputs := listInputs(b)
	require.Equal(t, common.expandedInputs, expandedInputs)
	mutant.RemoveAtoms(b)
	clippedInputs := listInputs(b)
	require.Equal(t, common.inputs, clippedInputs)
}
