package gblocks

import (
	r "reflect"
	"strconv"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/enum"
	"github.com/ionous/gblocks/mutant"
	"github.com/ionous/gblocks/option"
	"github.com/ionous/gblocks/tin"
)

type Maker struct {
	pairs    enum.Pairs
	mutables tin.Mutables // for buildItems, BlockStats
	types    []*tin.TypeInfo
}

func NewMaker(pairs enum.Pairs, mutations tin.Mutables, types []*tin.TypeInfo) *Maker {
	return &Maker{pairs, mutations, types}
}

func (m *Maker) GetPairs(n string) []enum.Pair {
	return m.pairs.GetPairs(n)
}

func (m *Maker) GetTerms(n string) (block.Limits, error) {
	return tin.LimitsOfOutput(m.types, n)
}

func (m *Maker) GetStatements(n string) (block.Limits, error) {
	return tin.LimitsOfNext(m.types, n)
}

func (m *Maker) GetTermsByType(t r.Type) block.Limits {
	return tin.LimitsOfType(m.types, t, tin.TermBlock)
}

func (m *Maker) GetStatementsByType(t r.Type) block.Limits {
	return tin.LimitsOfType(m.types, t, tin.MidBlock)
}

// submit all previously added blocks to a blockly-like project.
// mbs is the target for user generated mutation instance data ( where these blocks will write to if necessary )
// opts contains per-block type extra options.
func (m *Maker) RegisterBlocks(p block.Project, mbs mutant.MutatedBlocks, opts map[string]block.Dict) (err error) {
	// types is a slice of tin.TypeInfo
	for _, t := range m.types {
		if e := m.registerType(t, p, mbs, opts[t.Name]); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (m *Maker) registerType(t *tin.TypeInfo, p block.Project, mbs mutant.MutatedBlocks, opt block.Dict) (err error) {
	var mins mutant.InMutations
	if desc, e := m.makeDescByType(t, opt, &mins); e != nil {
		err = e
	} else {
		if mins.Mutates() {
			if e := mins.Preregister(t.Name, p); e != nil {
				err = e
			} else {
				mutatorName := mutant.MutatorName(t.Name)
				b := mutant.NewMutator(&mins, m, mbs)
				if e := p.RegisterMutator(mutatorName, b); e != nil {
					err = e
				} else {
					// add the registered mutation name to the block desc
					desc[option.Mutator] = mutatorName
				}
			}
		}
		if e := p.RegisterBlock(t.Name, desc); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// return a description of the named block in a format blockly can use for registration.
func (m *Maker) makeDesc(name string, out *mutant.InMutations) (ret block.Dict, err error) {
	if t, ok := tin.FindByName(m.types, name); !ok {
		err = errutil.New("unknown type", name)
	} else {
		ret, err = m.makeDescByType(t, nil, out)
	}
	return
}

//
func (m *Maker) makeDescByType(t *tin.TypeInfo, opt block.Dict, out *mutant.InMutations) (ret block.Dict, err error) {
	var db mutant.Atomizer = m
	if args, e := t.BuildItems("", db, m.mutables, out); e != nil {
		err = e
	} else {
		desc, hasContent := make(block.Dict), false
		if msg := args.Message(); len(msg) > 0 {
			block.Merge(desc, opt, option.Message(0), msg)
			hasContent = true
		}
		if list := args.List(); len(list) > 0 {
			block.Merge(desc, opt, option.Args(0), list)
			hasContent = true
		}
		if !hasContent {
			label := strings.Replace(t.Name, "_", " ", -1)
			block.Merge(desc, opt, option.Message(0), label)
		}
		if _, exists := desc[option.Tooltip]; !exists {
			tip := strings.Replace(t.Name, "_", " ", -1)
			block.Merge(desc, opt, option.Tooltip, tip)
		}
		// note: optional for jsonInit, mandatory for defineBlocksWithJsonArray
		block.Merge(desc, opt, option.Type, t.Name)
		//
		switch model := t.Model; model {
		case tin.MidBlock:
			block.Merge(desc, opt, option.Prev, t.Name)
			fallthrough

		case tin.TopBlock:
			if l := t.LimitsOfNext(m.types); l.Connects {
				block.Merge(desc, opt, option.Next, l.Check())
			}

		case tin.TermBlock:
			if l, e := t.LimitsOfOutput(m.types); e != nil {
				err = errutil.Append(err, e)
			} else if l.Connects {
				block.Merge(desc, opt, option.Output, l.Check())
			}

		default:
			panic(model.String())
		}
		if err == nil {
			// overwrite with block options
			t.VisitOptions(func(k, v string) {
				var i interface{}
				if n, e := strconv.ParseBool(v); e == nil {
					i = n
				} else if n, e := strconv.Atoi(v); e == nil {
					i = n
				} else {
					i = v
				}
				desc[k] = i
			})
			// overwrite with original options
			for k, v := range opt {
				desc[k] = v
			}
			ret = desc
		}
	}
	return
}
