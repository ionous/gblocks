package inspect

import (
	. "github.com/ionous/gblocks/gtest"
	"github.com/kr/pretty"
	"testing"
)

func TestPairs(t *testing.T) {
	var reg EnumPairs
	if pairs, e := reg.AddEnum(map[Enum]string{
		DefaultChoice:     "default",
		AlternativeChoice: "alt",
	}); e != nil {
		t.Fatal(e)
	} else if len(pairs) == 0 {
		t.Fatal("missing pairs")
	} else {
		expectedPairs := []EnumPair{
			{"default", "DefaultChoice"},
			{"alt", "AlternativeChoice"},
		}
		if v := pretty.Diff(pairs, expectedPairs); len(v) != 0 {
			t.Fatal(v)
		}
	}
}
