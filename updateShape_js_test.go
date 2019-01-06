package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
	r "reflect"
	"testing"
)

type ShapeTest struct {
	Input  *ShapeTest
	Mutant []interface{} `mutation:"TestMutation"`
	Field  string
}

type MutationEl struct {
	SubInput *ShapeTest
}

type MutationAlt struct {
	SubField string
}

type Verify struct {
	Name      InputName
	Type      InputType
	Mutations int
}

func reduce(b *Block) (ret []Verify) {
	mutations := func(in *Input) (ret int) {
		if m := in.Mutation(); m != nil {
			ret = m.TotalInputs()
		} else {
			ret = -1
		}
		return
	}
	for i := 0; i < b.NumInputs(); i++ {
		in := b.Input(i)
		ret = append(ret, Verify{
			Name:      in.Name,
			Type:      in.Type,
			Mutations: mutations(in),
		})
	}
	return
}

func TestShapeChangeless(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e)
		a1 := reduce(b)
		require.NoError(t, b.updateShape(ws))
		a2 := reduce(b)
		require.Equal(t, a1, a2)
		require.NoError(t, b.updateShape(ws))
		a3 := reduce(b)
		require.Equal(t, a2, a3)
	})
}

func TestShapeUpdate(t *testing.T) {
	subInput := Verify{"SUB_INPUT", InputValue, -1}
	//
	none := []Verify{{"INPUT", InputValue, -1}, {"MUTANT", DummyInput, 0}, {"", DummyInput, -1}}
	//
	v := [][]Verify{
		{{"INPUT", InputValue, -1}, {"MUTANT", DummyInput, 1}, subInput, {"", DummyInput, -1}},
		{{"INPUT", InputValue, -1}, {"MUTANT", DummyInput, 2}, subInput, subInput, {"", DummyInput, -1}},
		{{"INPUT", InputValue, -1}, {"MUTANT", DummyInput, 3}, subInput, subInput, subInput, {"", DummyInput, -1}},
		{{"INPUT", InputValue, -1}, {"MUTANT", DummyInput, 4}, subInput, subInput, subInput, subInput, {"", DummyInput, -1}},
	}

	testShape(t, func(ws *Workspace) {
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e)
		require.Equalf(t, none, reduce(b), "initially empty")
		// grow data which should grow the number of inputs.
		d := ws.GetDataById(b.Id).(*ShapeTest)
		for i := 0; i < 3; i++ {
			t.Log("adding element", i)
			d.Mutant = append(d.Mutant, &MutationEl{})
			require.NoErrorf(t, b.updateShape(ws), "element %d", i)
			require.Equalf(t, v[i], reduce(b), "element %d", i)
		}

		// reset data ( and inputs ) back to zero
		d.Mutant = nil
		require.NoError(t, b.updateShape(ws))
		require.Equalf(t, none, reduce(b), "ends empty")
	})
}

func TestShapeCreate(t *testing.T) {
	TheRegistry = Registry{}
	reg := &TheRegistry
	// field has unknown type Mutant gblocks.ShapeMutation
	require.NoError(t, RegisterMutation("TestMutation",
		nil, (*MutationElControl)(nil),
		(*MutationEl)(nil), (*MutationElControl)(nil),
		(*MutationAlt)(nil), (*MutationAltControl)(nil),
	), "register mutations")
	require.NoError(t, RegisterBlocks(nil,
		(*ShapeTest)(nil),
		(*MutationEl)(nil),
		(*MutationAlt)(nil),
		(*MutationElControl)(nil),
		(*MutationAltControl)(nil),
	), "register blocks")
	//
	var testShape = map[string]interface{}{
		"type":     TypeName("shape_test"),
		"message0": "%1 %2 %3",
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
		t.Fatal(v)
		t.Log(v)
	}
}

func testShape(t *testing.T, fn func(*Workspace)) {
	TheRegistry = Registry{}
	// field has unknown type Mutant gblocks.ShapeMutation
	require.NoError(t, RegisterMutation("TestMutation",
		nil, (*MutationElControl)(nil),
		(*MutationEl)(nil), (*MutationElControl)(nil),
		(*MutationAlt)(nil), (*MutationAltControl)(nil),
	), "register mutations")
	require.NoError(t, RegisterBlocks(nil,
		(*ShapeTest)(nil),
		(*MutationEl)(nil),
		(*MutationAlt)(nil),
		(*MutationElControl)(nil),
		(*MutationAltControl)(nil),
	), "register blocks")
	ws := NewBlankWorkspace()
	// replace timed event queue with direct event queue
	events := &Events{Object: js.Global.Get("Blockly").Get("Events")}
	events.Set("fire", js.MakeFunc(func(_ *js.Object, args []*js.Object) interface{} {
		events.TestFire(args[0])
		return nil
	}))
	fn(ws)
	ws.Dispose()
}
