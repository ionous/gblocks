package inspect

import (
	"github.com/ionous/gblocks/block"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
	r "reflect"
	"testing"
)

type CheckNext struct {
	NextStatement *CheckNext
}

type CheckStatement struct {
	PreviousStatement interface{}
	NextStatement     *CheckStatement
}

// TestConstraints - dependency pool generates constaints
func TestConstraints(t *testing.T) {
	var dp DependencyPool
	checkType, checkNext := r.TypeOf((*CheckStatement)(nil)), r.TypeOf((*CheckNext)(nil))
	if e := dp.AddTypes(checkType, checkNext); e != nil {
		t.Fatal("check statement", e)
	} else if f, ok := checkNext.Elem().FieldByName(NextStatement); !ok {
		t.Fatal("missing field")
	} else {
		constraints, ok := dp.GetConstraints(f.Type)
		require.True(t, ok, "get constraints")
		require.Equal(t, []block.Type{"check_next"}, constraints)
	}
}

// TestConstraintsAny - dependency pool generates constaints
func TestConstraintsAny(t *testing.T) {
	var dp DependencyPool
	checkType, checkNext := r.TypeOf((*CheckStatement)(nil)), r.TypeOf((*CheckNext)(nil))
	if e := dp.AddTypes(checkType, checkNext); e != nil {
		t.Fatal("check statement", e)
	} else if f, ok := checkType.Elem().FieldByName(PreviousStatement); !ok {
		t.Fatal("missing field")
	} else {
		constraints, ok := dp.GetConstraints(f.Type)
		require.True(t, ok, "get constraints")
		require.Equal(t, []block.Type(nil), constraints)
	}
}

// TestConstraintDesc - ensure the type pool can produce next and prev links with appropriate constraints.
func TestConstraintDesc(t *testing.T) {
	var tp TypePool
	checkType, checkNext := r.TypeOf((*CheckStatement)(nil)), r.TypeOf((*CheckNext)(nil))
	if e := tp.AddTypes(checkType, checkNext); e != nil {
		t.Fatal(e)
	} else if desc, e := tp.BuildDesc(checkType, nil); e != nil {
		t.Fatal(e)
	} else {
		t.Log(pretty.Sprint(desc))
		expected := block.Dict{
			"message0":          "check statement", // has no fields, so uses its own name.
			"previousStatement": []block.Type(nil),
			"nextStatement":     []block.Type{"check_statement"},
			"type":              block.Type("check_statement"),
		}
		if v := pretty.Diff(desc, expected); len(v) != 0 {
			t.Fatal(v)
		}
	}
}
