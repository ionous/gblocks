package inspect

import (
	"github.com/kr/pretty"
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

func TestPairs(t *testing.T) {
	var reg EnumPairs
	if pairs, e := reg.RegisterEnum(map[Enum]string{
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
