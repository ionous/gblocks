package gblocks

import (
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
	r "reflect"
	"strconv"
	"testing"
)

type ShapeTest struct {
	Input  *ShapeTest
	Mutant TestMutation
	Field  string
}

type TestMutation struct {
	NextStatement NextAtom
}

// Output - implement a generic output
func (n *ShapeTest) Output() *ShapeTest {
	return n
}

// NextAtom
type NextAtom interface {
	NextAtom() NextAtom
}

type AtomTest struct {
	AtomInput     *ShapeTest
	NextStatement NextAtom
}

func (a *AtomTest) NextAtom() NextAtom { return a.NextStatement }

type AtomAltTest struct {
	AtomField     string
	NextStatement NextAtom
}

func (a *AtomAltTest) NextAtom() NextAtom { return a.NextStatement }

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
	//
	require.NoError(t,
		reg.RegisterMutation((*TestMutation)(nil)),
		"register mutations")
	//
	require.NoError(t,
		reg.RegisterBlock((*ShapeTest)(nil), nil),
		"register blocks")
	//
	var testShape = Dict{
		"type":     TypeName("shape_test"),
		"message0": "%1 %2 %3",
		"output":   TypeName("shape_test"),
		"args0": []Dict{
			{
				"name":  "INPUT",
				"type":  "input_value",
				"check": []TypeName{"shape_test"},
			},
			{
				"name":     "MUTANT",
				"type":     "input_dummy",
				"mutation": TypeName("test_mutation"),
			},
			{
				"name": "FIELD",
				"type": "field_input",
				"text": "Field",
			},
		},
	}
	opt := make(Dict)
	reg.buildBlockDesc(r.TypeOf((*ShapeTest)(nil)).Elem(), opt)
	if v := pretty.Diff(opt, testShape); len(v) != 0 {
		t.Log(pretty.Sprint(opt))
		t.Fatal(v)
	}
}

func testShape(t *testing.T, fn func(*Workspace, *Registry)) {
	reg := new(Registry)
	require.NoError(t,
		reg.RegisterMutation((*TestMutation)(nil),
			Mutation{"atom", (*AtomTest)(nil)},
			Mutation{"alt", (*AtomAltTest)(nil)},
		), "register mutations")
	//
	require.NoError(t, reg.RegisterBlocks(nil,
		(*ShapeTest)(nil),
		(*AtomTest)(nil),
		(*AtomAltTest)(nil),
	), "register blocks")
	ws := NewBlankWorkspace(false)
	ws.idGen = &orderedGenerator{name: "main"}
	fn(ws, reg)
	ws.Dispose()
}
