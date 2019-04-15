package mutant

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/option"
)

type Quarks interface {
	Quarks(paletteOnly bool) (Quark, bool)
}

// find a quark with the passed full or short name
func FindQuark(quarks Quarks, name string) (ret Quark, okay bool) {
	for q, ok := quarks.Quarks(false); ok; q, ok = q.NextQuark() {
		if name == q.Name() || name == q.BlockType() {
			ret, okay = q, true
			break
		}
	}
	return
}

// return a list of blocks used in a palette; exludes the fixed block if any.
func PaletteQuarks(quarks Quarks) (ret []string) {
	for q, ok := quarks.Quarks(true); ok; q, ok = q.NextQuark() {
		ret = append(ret, q.BlockType())
	}
	return
}

// make blockly compatible description of the quark's mui block
func DescribeQuark(q Quark) block.Dict {
	out := block.Dict{
		option.Type:       q.BlockType(),
		option.Message(0): q.Label(),
		option.Prev:       q.Name(),
	}
	if next := q.LimitsOfNext(); next.Connects {
		out[option.Next] = next.Check()
	}
	return out
}

func RegisterQuarks(p block.Project, qs Quarks) (err error) {
	for q, ok := qs.Quarks(false); ok; q, ok = q.NextQuark() {
		name, desc := q.BlockType(), DescribeQuark(q)
		if e := p.RegisterBlock(name, desc); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
