package gblocks

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/inspect"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
	r "reflect"
	"testing"
)

type CheckNext struct {
	NextStatement *CheckNext
}

type CheckPrev struct {
	PreviousStatement *CheckPrev
}

type CheckStatement struct {
	NextStatement, PreviousStatement *CheckStatement
}

func TestCheckStatement(t *testing.T) {
	var reg Registry
	checkType := r.TypeOf((*CheckStatement)(nil))
	if desc, e := reg.testRegister(checkType); e != nil {
		t.Fatal(e)
	} else {
		t.Log(pretty.Sprint(desc))
		expected := block.Dict{
			"message0":          "check statement", // has no fields, so uses its own name.
			"previousStatement": []block.Type{"check_statement"},
			"nextStatement":     []block.Type{"check_statement"},
			"type":              block.Type("check_statement"),
		}
		v := pretty.Diff(desc, expected)
		if len(v) != 0 {
			t.Fatal(v)
		}
	}
}

func TestCheckNext(t *testing.T) {
	var reg Registry
	checkType, ptrType := r.TypeOf((*CheckStatement)(nil)), r.TypeOf((*CheckNext)(nil))
	if _, e := reg.testRegister(checkType); e != nil {
		t.Fatal("check statement", e)
	} else if _, e := reg.testRegister(ptrType); e != nil {
		t.Fatal("next statement", e)
	} else if f, ok := ptrType.Elem().FieldByName(inspect.NextStatement); !ok {
		t.Fatal("missing field")
	} else {
		constraints, ok := reg.GetConstraints(f.Type)
		require.True(t, ok, "get constraints")
		require.Equal(t, constraints, []block.Type{"check_next"})
	}
}

func TestCheckPrev(t *testing.T) {
	var reg Registry
	checkType, ptrType := r.TypeOf((*CheckStatement)(nil)), r.TypeOf((*CheckPrev)(nil))
	if _, e := reg.testRegister(checkType); e != nil {
		t.Fatal("check statement", e)
	} else if _, e := reg.testRegister(ptrType); e != nil {
		t.Fatal("prev statement", e)
	} else if f, ok := ptrType.Elem().FieldByName(inspect.PreviousStatement); !ok {
		t.Fatal("missing field")
	} else {
		constraints, ok := reg.GetConstraints(f.Type)
		require.True(t, ok, "get constraints")
		require.Equal(t, constraints, []block.Type{"check_prev"})
	}
}
