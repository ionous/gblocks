package mutant

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/option"
)

// descriptions of the mutable inputs for a given type of block.
type BlockMutations struct {
	Inputs    []string            // ordered input names (b/c golang maps are unordered)
	Mutations map[string]Mutation // input name to mutation info interface
}

func (m *BlockMutations) AddMutation(inputName string, arch Mutation) {
	if m.Mutations == nil {
		m.Mutations = make(map[string]Mutation)
	}
	m.Inputs = append(m.Inputs, inputName)
	m.Mutations[inputName] = arch
}

// does this block type have any mutations?
func (m *BlockMutations) Mutates() bool {
	return len(m.Inputs) > 0
}

// return the mutation description for the passed input name.
func (m *BlockMutations) GetMutation(n string) (ret Mutation, okay bool) {
	ret, okay = m.Mutations[n]
	return
}

// visit all quarks defined by all input mutations
func (m *BlockMutations) Quarks(paletteOnly bool) (ret Quark, okay bool) {
	it := muteIt{m, -1, nil, paletteOnly}
	return it.advance()
}

// produce a description of the mui container.
// each input in the container represents a mutable input in the workspace block.
func (m *BlockMutations) DescribeContainer(containerName string) (ret block.Dict) {
	var args block.Args
	for _, in := range m.Inputs {
		arch := m.Mutations[in] // Mutation
		l := arch.Limits()
		arg := block.Dict{
			option.Name: in,
			option.Type: block.StatementInput,
		}
		// DescribeQuark sets a prev of the quark name.
		// we need to limit the input to those acceptable types
		// we're asking for the "mutation" generically
		// normally we'd what "all available types" to be "no limit"
		// but that doesnt work well here -- where we are mixing these multiple types together
		if !l.IsUnlimited() {
			arg[option.Check] = l.Check()
		}
		args.AddArg(arg)
	}
	return block.Dict{
		option.Type:       containerName,
		option.Message(0): args.Message(),
		option.Args(0):    args.List(),
	}
}

// register the mui container and all of its mui blocks with blockly
func (m *BlockMutations) Preregister(blockType string, p block.Project) (err error) {
	containerName := ContainerName(blockType)
	desc := m.DescribeContainer(containerName)
	if e := p.RegisterBlock(containerName, desc); e != nil {
		err = e
	} else if e := RegisterQuarks(p, m); e != nil {
		err = e
	}
	return
}
