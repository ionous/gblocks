package mutant

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
	"github.com/ionous/gblocks/option"
)

// acts as an ordered map of input name to input mutation.
// a blockly mutator consists of all the input mutations for a given block.
type InMutations struct {
	Inputs    []string              // ordered input names (b/c go maps arent ordered)
	Mutations map[string]InMutation // input name to mutation data
}

func (m *InMutations) Mutates() bool {
	return len(m.Inputs) > 0
}

func (m *InMutations) GetMutation(n string) (ret InMutation, okay bool) {
	ret, okay = m.Mutations[n]
	return
}

// visit all quarks defined by all input mutations
func (m *InMutations) Quarks(paletteOnly bool) (ret Quark, okay bool) {
	it := muteIt{m, -1, nil, paletteOnly}
	return it.advance()
}

// produce a description of the mui container.
// each input in the container represents a mutable input in the workspace block.
func (m *InMutations) DescribeContainer(containerName string) (ret block.Dict) {
	var args block.Args
	for _, in := range m.Inputs {
		min := m.Mutations[in] // InMutation
		l := min.Limits()
		arg := block.Dict{
			option.Name: in,
			option.Type: block.StatementInput,
		}
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
func (m *InMutations) Preregister(blockType string, p block.Project) (err error) {
	containerName := ContainerName(blockType)
	desc := m.DescribeContainer(containerName)
	if e := p.RegisterBlock(containerName, desc); e != nil {
		err = e
	} else if e := RegisterQuarks(p, m); e != nil {
		err = e
	}
	return
}

// aka "decompose" -- Populate the mutator popup with this block's components.
func (m *InMutations) CreateMui(mui block.Workspace, b block.Shape, inputs MutableInputs) (ret block.Shape, err error) {
	containerName := ContainerName(b.BlockType())
	if container, e := mui.NewBlock(containerName); e != nil {
		err = e
	} else {
		l := muiBuilder{m, b.BlockId(), container, inputs}
		if e := l.fillContainer(); e != nil {
			container.Dispose()
			err = e
		} else {
			ret = container
		}
	}
	return
}

// aka. compose.
// adds new inputs to target, returns the atoms for those inputs
func (m *InMutations) DistillMui(target, muiContainer block.Shape, db Atomizer, cs SavedConnections) (ret MutableInputs, err error) {
	// remove all the dynamic inputs from the blocks; we're about to recreate/recompose them.
	// note: the connections for those inputs have already been saved in cs.
	RemoveAtoms(target)
	mp := muiParser{m, target, db}
	if atoms, e := mp.expandInputs(muiContainer); e != nil {
		err = errutil.New("fill from mui", e)
	} else {
		RestoreConnections(target, cs) // no return
		ret = atoms
	}
	return
}

// aka. domToMutation
func (m *InMutations) LoadMutation(b block.Shape, items Atomizer, mutationEls dom.BlockMutation) (ret MutableInputs, err error) {
	dp := domParser{m, items, b, make(MutableInputs)}
	if e := dp.parseDom(&mutationEls); e != nil {
		err = errutil.New("load mutation", e)
	} else {
		ret = dp.inputs
	}
	return
}

// serialize the mutation to xml friendly data
// we use a structured intermediary to simplify testing in go.
func (m *InMutations) SaveMutation(inputs MutableInputs) (ret dom.BlockMutation) {
	for _, inputName := range m.Inputs {
		if atoms, ok := inputs[inputName]; ok {
			// if there are atoms, create a node for the data.
			if numAtoms := len(atoms); numAtoms > 0 {
				ret.Append(&dom.Mutation{inputName, dom.Atoms{atoms}})
			} else {
				// no explict atoms? there might be an implicit first block
				in := m.Mutations[inputName]
				if first, ok := in.FirstBlock(); ok {
					atoms := []string{first.Name()}
					ret.Append(&dom.Mutation{inputName, dom.Atoms{atoms}})
				}
			}
		}
	}
	return
}
