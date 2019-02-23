package gblocks

import (
	. "github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/gtest"
	"github.com/ionous/gblocks/inspect"
	"github.com/kr/pretty"
	r "reflect"
	"testing"
)

func TestEnumLabels(t *testing.T) {
	var reg Registry
	if _, e := reg.RegisterEnum(map[Enum]string{
		DefaultChoice:     "default",
		AlternativeChoice: "alt",
	}); e != nil {
		t.Fatal(e)
	} else {
		if blockDesc, e := reg.testRegister(r.TypeOf((*EnumStatement)(nil))); e != nil {
			t.Fatal(e)
		} else {
			t.Log(pretty.Sprint(blockDesc))
			expected := block.Dict{
				"message0": "%1",
				"args0": []block.Dict{
					{
						"name": "ENUM",
						"type": "field_dropdown",
						"options": []inspect.EnumPair{
							{"default", "DefaultChoice"},
							{"alt", "AlternativeChoice"},
						},
					},
				},
				"type": block.Type("enum_statement"),
			}
			v := pretty.Diff(blockDesc, expected)
			if len(v) != 0 {
				t.Fatal(v)
			}
		}
	}
}
