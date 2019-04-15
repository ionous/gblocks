package mock

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
	"github.com/ionous/gblocks/option"
)

type MockQuark struct {
	mutation *MockMutation
	i        int
}

// type name without including owner mutation
func (q *MockQuark) Name() string {
	name := q.mutation.quarks[q.i]
	return name
}

// type name scoped to the owner mutation
func (q *MockQuark) BlockType() string {
	name := q.mutation.quarks[q.i]
	return block.Scope("mui", q.mutation.name, name)
}

func (q *MockQuark) Label() string {
	name := q.mutation.quarks[q.i]
	return name
}

// make blockly compatible description of the quark's mui block
func (q *MockQuark) LimitsOfNext() block.Limits {
	return block.MakeUnlimited()
}

// expand the quark into an atom ( a bundle of fields and inputs )
func (q *MockQuark) Atomize(scope string, db mutant.Atomizer) (block.Args, error) {
	var args block.Args
	atoms := db.(*MockDatabase).Atoms[q.Name()]
	for _, a := range atoms {
		args.AddArg(block.Dict{
			// see also gblocks.BuildItems
			option.Name: block.Scope(scope, a.Name),
			option.Type: a.Type,
		})
	}
	return args, nil
}

func (q *MockQuark) NextQuark() (ret mutant.Quark, okay bool) {
	if n := q.i + 1; n < len(q.mutation.quarks) {
		ret, okay = &MockQuark{q.mutation, n}, true
	}
	return
}
