package mutant

import (
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
)

// <mutation>
//   <pin name="MUTANT">
//     <atom type="decor_mutation"/>
//     <atom type="decor_mutation"/>
//   </pin>
// </mutation>
type domParser struct {
	mins   *InMutations  // description of mutatable inputs
	db     Atomizer      // to expand atoms into inputs
	block  block.Shape   // target block
	inputs MutableInputs // info on expanded inputs
}

type atomParser struct {
	*domParser
	min InMutation
	*injector
}

func (dp *domParser) parseDom(ms *dom.BlockMutation) (err error) {
	for _, el := range ms.Inputs {
		if e := dp.parseInput(el); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// read <pin>
func (dp *domParser) parseInput(inputEl *dom.Mutation) (err error) {
	name, atoms := inputEl.Input, inputEl.Atoms
	// watch for empty mutations ( where no input data exists )
	if len(name) > 0 {
		if min, ok := dp.mins.GetMutation(name); !ok {
			err = errutil.New("input not mutable", name)
		} else if j, e := newInjector(dp.block, name); e != nil {
			err = e
		} else {
			ap := atomParser{dp, min, j}
			if atoms, e := ap.parseAtoms(atoms.Types); e != nil {
				err = errutil.New("parsing", name, e)
			} else {
				dp.inputs[name] = atoms
				ap.finalizeInputs()
			}
		}
	}
	return
}

// parse children of <pin>
func (ap *atomParser) parseAtoms(atoms []string) (ret []string, err error) {
	for i, el := range atoms {
		if atom, e := ap.parseAtom(el, i); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = append(ret, atom)
		}
	}
	return
}

// parse an individual <atom>
func (ap *atomParser) parseAtom(quark string, atomNum int) (ret string, err error) {
	if q, ok := FindQuark(ap.min, quark); !ok {
		err = errutil.New("quark not found", quark)
	} else {
		wsBlockId, inputName := ap.block.BlockId(), ap.target().InputName()
		atomScope := block.Scope("a", wsBlockId, inputName, strconv.Itoa(atomNum))
		if args, e := q.Atomize(atomScope, ap.db); e != nil {
			err = e
		} else if e := ap.inject(args); e != nil {
			err = e
		} else {
			ret = q.Name()
		}
	}
	return
}
