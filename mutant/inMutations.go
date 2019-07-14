package mutant

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
	"github.com/ionous/gblocks/option"
)

// acts as an ordered map of input name to input mutation.
// a blockly mutator consists of all the input mutations for a given block.
type BlockMutations struct {
	Inputs    []string            // ordered input names (b/c go maps arent ordered)
	Mutations map[string]Mutation // input name to mutation data
}

func (m *BlockMutations) Mutates() bool {
	return len(m.Inputs) > 0
}

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
		min := m.Mutations[in] // InMutation
		l := min.Limits()
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

// aka "decompose" -- Populate the mutator popup with this block's components.
func (m *BlockMutations) CreateMui(mui block.Workspace, b block.Shape, inputs AtomizedInputs) (ret block.Shape, err error) {
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

// aka. compose -- turn the mui into new workspace inputs
// adds new inputs to target, returns the atoms for those inputs
func (m *BlockMutations) DistillMui(target, muiContainer block.Shape, db Atomizer, cs SavedConnections) (ret AtomizedInputs, err error) {
	// remove all the dynamic inputs from the blocks; we're about to recreate/recompose them.
	// note: the connections for those inputs have already been saved in cs.
	RemoveAtoms(target)
	mp := muiParser{m, target, db}
	if atoms, e := mp.expandInputs(muiContainer); e != nil {
		err = errutil.New("DistillMui()", e)
	} else {
		RestoreConnections(target, cs) // no return
		ret = atoms
	}
	return
}

// aka. domToMutation
func (m *BlockMutations) LoadMutation(b block.Shape, items Atomizer, mutationEls dom.BlockMutation) (ret AtomizedInputs, err error) {
	dp := domParser{m, items, b, MakeAtomizedInputs()}
	if e := dp.parseDom(&mutationEls); e != nil {
		err = errutil.New("LoadMutation()", e)
	} else {
		ret = dp.inputs
	}
	return
}

// serialize the mutation to xml friendly data
// we use a structured intermediary to simplify testing in go.
func (m *BlockMutations) SaveMutation(inputs AtomizedInputs) (ret dom.BlockMutation) {
	for _, inputName := range m.Inputs {
		if atoms, ok := inputs.GetAtomsForInput(inputName); ok {
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
