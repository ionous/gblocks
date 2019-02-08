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
}

func TestEnumLabels(t *testing.T) {
	var reg Registry
	reg.enums.registerEnum(map[Enum]string{
		DefaultChoice:     "default",
		AlternativeChoice: "alt",
	})
	desc := make(Dict)
	reg.buildBlockDesc(r.TypeOf((*EnumStatement)(nil)).Elem(), desc)
	t.Log(pretty.Sprint(desc))
	expected := Dict{
		"message0": "%1",
		"args0": []Dict{
			{
				"name": "ENUM",
				"type": "field_dropdown",
				"options": []EnumPair{
					{"default", "DefaultChoice"},
					{"alt", "AlternativeChoice"},
				},
			},
		},
		"type": TypeName("enum_statement"),
	}
	v := pretty.Diff(desc, expected)
	if len(v) != 0 {
		t.Fatal(v)
	}
}
