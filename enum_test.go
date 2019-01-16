package gblocks

import (
	"github.com/kr/pretty"
	r "reflect"
	"testing"
)

type Enum int

const (
	DefaultChoice Enum = iota
	AlternativeChoice
)

func (i Enum) String() (ret string) {
	switch i {
	case DefaultChoice:
		ret = "DefaultChoice"
	case AlternativeChoice:
		ret = "AlternativeChoice"
	}
	return
}

type EnumStatement struct {
	Enum
	PreviousStatement, NextStatement *EnumStatement
}

func TestEnumLabels(t *testing.T) {
	var reg Registry
	reg.registerEnum(map[Enum]string{
		DefaultChoice:     "default",
		AlternativeChoice: "alt",
	})
	opt := make(map[string]interface{})
	reg.initJson(r.TypeOf((*EnumStatement)(nil)).Elem(), opt)
	expected := map[string]interface{}{
		"message0": "%1",
		"args0": []Options{
			{
				"name": "ENUM",
				"type": "field_dropdown",
				"options": []stringPair{
					{"default", "DefaultChoice"},
					{"alt", "AlternativeChoice"},
				},
			},
		},
		"previousStatement": TypeName("enum_statement"),
		"nextStatement":     TypeName("enum_statement"),
		"type":              TypeName("enum_statement"),
	}
	v := pretty.Diff(opt, expected)
	if len(v) != 0 {
		t.Fatal(v)
		t.Log(v)
	}
}
