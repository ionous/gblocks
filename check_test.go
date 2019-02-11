package gblocks

import (
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

func TestCheckNext(t *testing.T) {
	types := make(RegisteredTypes)
	structType := r.TypeOf((*CheckNext)(nil)).Elem()
	if ok := types.RegisterType(structType); !ok {
		t.Fatal("couldnt register type")
	} else {
		check, e := types.CheckField(structType, NextField)
		require.NoError(t, e, "check field")
		constraints, ok := check.GetConstraints()
		require.True(t, ok, "get constraints")
		require.Equal(t, constraints, []TypeName{"check_next"})
	}
}

func TestCheckPrev(t *testing.T) {
	types := make(RegisteredTypes)
	structType := r.TypeOf((*CheckPrev)(nil)).Elem()
	if ok := types.RegisterType(structType); !ok {
		t.Fatal("couldnt register type")
	} else {
		check, e := types.CheckField(structType, PreviousField)
		require.NoError(t, e, "check field")
		constraints, ok := check.GetConstraints()
		require.True(t, ok, "get constraints")
		require.Equal(t, constraints, []TypeName{"check_prev"})
	}
}

func TestCheckDesc(t *testing.T) {
	types := make(RegisteredTypes)
	structType := r.TypeOf((*CheckStatement)(nil)).Elem()
	if ok := types.RegisterType(structType); !ok {
		t.Fatal("couldnt register type")
	} else {
		reg := Registry{types: types}
		desc := make(Dict)
		if _, e := reg.buildBlockDesc(r.TypeOf((*CheckStatement)(nil)).Elem(), desc); e != nil {
			t.Fatal("couldnt describe block", e)
		} else {
			t.Log(pretty.Sprint(desc))
			expected := Dict{
				"message0":          "check statement", // has no fields, so uses its own name.
				"previousStatement": []TypeName{"check_statement"},
				"nextStatement":     []TypeName{"check_statement"},
				"type":              TypeName("check_statement"),
			}
			v := pretty.Diff(desc, expected)
			if len(v) != 0 {
				t.Fatal(v)
			}
		}
	}
}
