package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
	r "reflect"
	"strconv"
	"testing"
)

type ShapeTest struct {
	Input  *ShapeTest
	Mutant []interface{} `mutation:"TestMutation"`
	Field  string
}

// Output - implement a generic output
func (n *ShapeTest) Output() *ShapeTest {
	return n
}

type AtomTest struct {
	AtomInput *ShapeTest
}

type AtomAltTest struct {
	AtomField string
}

type orderedGenerator struct {
	name string
	i    int
}

func (o *orderedGenerator) NewId() string {
	o.i++
	return o.name + strconv.Itoa(o.i)
}

func TestShapeCreate(t *testing.T) {
	var reg Registry
	// field has unknown type Mutant gblocks.ShapeMutation
	require.NoError(t, reg.RegisterMutation("TestMutation",
		nil, (*MutationElControl)(nil),
		(*AtomTest)(nil), (*MutationElControl)(nil),
		(*AtomAltTest)(nil), (*MutationAltControl)(nil),
	), "register mutations")
	require.NoError(t, reg.RegisterBlocks(nil,
		(*ShapeTest)(nil),
		(*AtomTest)(nil),
		(*AtomAltTest)(nil),
		(*MutationElControl)(nil),
		(*MutationAltControl)(nil),
	), "register blocks")
	//
	var testShape = map[string]interface{}{
		"type":     TypeName("shape_test"),
		"message0": "%1 %2 %3",
		"output":   TypeName("shape_test"),
		"args0": []Options{
			{
				"name":  "INPUT",
				"type":  "input_value",
				"check": TypeName("shape_test"),
			},
			{
				"mutation": "TestMutation",
				"name":     "MUTANT",
				"type":     "input_dummy",
			},
			{
				"name": "FIELD",
				"type": "field_input",
				"text": "Field",
			},
		},
	}
	opt := make(map[string]interface{})
	reg.initJson(r.TypeOf((*ShapeTest)(nil)).Elem(), opt)
	if v := pretty.Diff(opt, testShape); len(v) != 0 {
		t.Log(pretty.Sprint(opt))
		t.Fatal(v)
	}
}

func testShape(t *testing.T, fn func(*Workspace, *Registry)) {
	reg := new(Registry)
	// field has unknown type Mutant gblocks.ShapeMutation
	require.NoError(t, reg.RegisterMutation("TestMutation",
		nil, (*MutationElControl)(nil),
		(*AtomTest)(nil), (*MutationElControl)(nil),
		(*AtomAltTest)(nil), (*MutationAltControl)(nil),
	), "register mutations")
	//
	require.NoError(t, reg.RegisterBlocks(nil,
		(*ShapeTest)(nil),
		(*AtomTest)(nil),
		(*AtomAltTest)(nil),
		(*MutationElControl)(nil),
		(*MutationAltControl)(nil),
	), "register blocks")
	ws := NewBlankWorkspace(false)
	ws.idGen = &orderedGenerator{name: "main"}
	// replace timed event queue with direct event queue
	events := GetEvents()
	events.Set("fire", js.MakeFunc(func(_ *js.Object, args []*js.Object) interface{} {
		events.TestFire(args[0])
		return nil
	}))
	fn(ws, reg)
	ws.Dispose()
}
