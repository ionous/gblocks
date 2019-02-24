package inspect

import (
	"github.com/ionous/gblocks/block"
	. "github.com/ionous/gblocks/gtest"
	"github.com/kr/pretty"
	r "reflect"
	"testing"
)

func TestAddEnum(t *testing.T) {
	var tp TypePool
	if _, e := tp.AddEnum(map[Enum]string{
		DefaultChoice:     "default",
		AlternativeChoice: "alt",
	}); e != nil {
		t.Fatal(e)
	} else {
		enumType := r.TypeOf((*EnumStatement)(nil))
		if e := tp.AddType(enumType); e != nil {
			t.Fatal(e)
		} else if blockDesc, e := tp.BuildDesc(enumType, nil); e != nil {
			t.Fatal(e)
		} else {
			t.Log(pretty.Sprint(blockDesc))
			expected := block.Dict{
				"message0": "%1",
				"args0": []block.Dict{
					{
						"name": "ENUM",
						"type": "field_dropdown",
						"options": []EnumPair{
							{"default", "DefaultChoice"},
							{"alt", "AlternativeChoice"},
						},
					},
				},
				"type": block.Type("enum_statement"),
			}
			if v := pretty.Diff(blockDesc, expected); len(v) != 0 {
				t.Fatal(v)
			}
		}
	}
}
