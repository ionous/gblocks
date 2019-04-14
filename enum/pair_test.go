package enum

import (
	"testing"

	. "github.com/ionous/gblocks/test"
	"github.com/kr/pretty"
)

func TestPairs(t *testing.T) {
	var reg Pairs
	if pairs, e := reg.AddEnum(map[Enum]string{
		DefaultChoice:     "default",
		AlternativeChoice: "alt",
	}); e != nil {
		t.Fatal(e)
	} else if len(pairs) == 0 {
		t.Fatal("missing pairs")
	} else {
		expectedPairs := []Pair{
			{"default", "DefaultChoice"},
			{"alt", "AlternativeChoice"},
		}
		if v := pretty.Diff(pairs, expectedPairs); len(v) != 0 {
			t.Fatal(v)
		}
	}
}
