package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/stretchr/testify/require"
	r "reflect"
	"testing"
)

type ShapeTest struct {
	Input  *ShapeTest
	Mutant ShapeMutation
	Field  string
}

type MutationEl struct {
	SubInput *ShapeTest
}

type MutationAlt struct {
	SubField string
}

type ShapeMutation struct {
	els []interface{}
}

func (m *ShapeMutation) Elements() r.Value {
	return r.ValueOf(m.els)
}

// MutationForType - given the passed data type; what block type is needed
func (m *ShapeMutation) MutationForType(t r.Type) (ret r.Type) {
	switch t {
	case nil:
		ret = r.TypeOf((*MutationElStart)(nil)).Elem()

	case r.TypeOf((*MutationEl)(nil)).Elem():
		ret = r.TypeOf((*MutationElControl)(nil)).Elem()

	case r.TypeOf((*MutationAlt)(nil)).Elem():
		ret = r.TypeOf((*MutationAltControl)(nil)).Elem()
	}
	return
}

type Verify struct {
	Name      string
	Type      InputType
	Mutations int
}

func reduce(b *Block) (ret []Verify) {
	for i := 0; i < b.NumInputs(); i++ {
		in := b.Input(i)
		ret = append(ret, Verify{
			Name:      in.Name,
			Type:      in.Type,
			Mutations: in.mutations,
		})
	}
	return
}

func TestShapeChangeless(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		b := ws.NewBlock((*ShapeTest)(nil))
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
	subInput := Verify{"SUB_INPUT", InputValue, 0}
	none := []Verify{{"INPUT", InputValue, 0}, {"MUTANT", DummyInput, -1}, {"", DummyInput, 0}}
	v := [][]Verify{
		{{"INPUT", InputValue, 0}, {"MUTANT", DummyInput, 1}, subInput, {"", DummyInput, 0}},
		{{"INPUT", InputValue, 0}, {"MUTANT", DummyInput, 2}, subInput, subInput, {"", DummyInput, 0}},
		{{"INPUT", InputValue, 0}, {"MUTANT", DummyInput, 3}, subInput, subInput, subInput, {"", DummyInput, 0}},
		{{"INPUT", InputValue, 0}, {"MUTANT", DummyInput, 4}, subInput, subInput, subInput, subInput, {"", DummyInput, 0}},
	}

	testShape(t, func(ws *Workspace) {
		b := ws.NewBlock((*ShapeTest)(nil))
		require.Equalf(t, none, reduce(b), "initially empty")
		// grow data which should grow the number of inputs.
		d := ws.GetDataById(b.Id).(*ShapeTest)
		for i := 0; i < 3; i++ {
			t.Log("adding element", i)
			d.Mutant.els = append(d.Mutant.els, &MutationEl{})
			require.NoErrorf(t, b.updateShape(ws), "element %d", i)
			require.Equalf(t, v[i], reduce(b), "element %d", i)
		}

		// reset data ( and inputs ) back to zero
		d.Mutant.els = nil
		require.NoError(t, b.updateShape(ws))
		require.Equalf(t, none, reduce(b), "ends empty")
	})
}

func testShape(t *testing.T, fn func(*Workspace)) {
	var reg Registry
	// field has unknown type Mutant gblocks.ShapeMutation
	require.NoError(t, reg.RegisterBlocks(nil,
		(*ShapeTest)(nil),
		(*MutationEl)(nil),
		(*MutationAlt)(nil),
		(*MutationElStart)(nil),
		(*MutationElControl)(nil),
		(*MutationAltControl)(nil),
	))
	ws := NewBlankWorkspace(&reg)
	// replace timed event queue with direct event queue
	events := &Events{Object: js.Global.Get("Blockly").Get("Events")}
	events.Set("fire", js.MakeFunc(func(_ *js.Object, args []*js.Object) interface{} {
		events.TestFire(args[0])
		return nil
	}))
	//
	fn(ws)
	ws.Dispose()
}
