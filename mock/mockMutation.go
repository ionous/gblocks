package mock

import (
	"strings"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
)

// a single input's mutation
type MockMutation struct {
	name   string // name of the mutation; mocked as its corresponding input name
	limits block.Limits
	quarks []string
}

// create unique mutations for every input
// based on one common set of quarks.
func NewMutations(inputs []string, quarks []string) *mutant.BlockMutations {
	mutations := make(map[string]mutant.Mutation)
	var names []string
	for _, input := range inputs {
		parts := strings.Split(input, ":")
		inputName := parts[0]
		mutation := &MockMutation{
			inputName,
			block.MakeLimits([]string{quarks[0]}),
			quarks,
		}
		mutations[inputName] = mutation
		names = append(names, inputName)
	}
	return &mutant.BlockMutations{names, mutations}
}

func (m *MockMutation) String() string {
	return "mock:" + m.name
}

func (m *MockMutation) Name() string {
	return m.name
}

// returns quarks which can appear at the top of the stack of the mui ui.
func (m *MockMutation) Limits() block.Limits {
	return m.limits
}

// no fixed block
func (m *MockMutation) FirstBlock() (mutant.Quark, bool) {
	return nil, false
}

func (m *MockMutation) Quarks(paletteOnly bool) (ret mutant.Quark, okay bool) {
	if len(m.quarks) > 0 {
		ret, okay = &MockQuark{m, 0}, true
	}
	return
}
