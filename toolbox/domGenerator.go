package toolbox

import (
	r "reflect" // for inspecting go-lang values

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
	"github.com/ionous/gblocks/option"
	"github.com/ionous/gblocks/pascal"
	"github.com/ionous/gblocks/tin"
)

// interface to help generate toolbox dom from golang structs
type Events interface {
	OnBlock(*tin.TypeInfo)
	OnError(error)
}

type UniqueNames interface {
	GenerateUniqueName() string
}

type domGenerator struct {
	events Events
	names  UniqueNames
}

// v should be a kind of struct
// returns a <block> (or <shadow>)
func (gen *domGenerator) genBlock(blockVal r.Value, shadowing Shadowing) *dom.Block {
	blockType := blockVal.Type()
	blockName := pascal.ToUnderscore(blockType.Name())
	block := &dom.Block{Type: blockName, Shadow: shadowing == IsShadow}
	g := blockGen{block: block, domGenerator: gen, shadowing: shadowing}
	g.fieldsOf(blockVal, blockType)
	return g.block
}

type blockGen struct {
	block *dom.Block
	*domGenerator
	shadowing Shadowing
	mutating  *mutatorGen
}

type mutatorGen struct {
	atoms     []*dom.Atom // generated atoms ( for generating input names )
	atomIndex int         // <0 handling the mutator struct, not the atoms
}

func (b *mutatorGen) currentAtom() (ret *dom.Atom, okay bool) {
	if b.atomIndex >= 0 && b.atomIndex < len(b.atoms) {
		ret, okay = b.atoms[b.atomIndex], true
	}
	return
}

func (b *mutatorGen) nextAtom() *mutatorGen {
	return &mutatorGen{b.atoms, b.atomIndex + 1}
}

func (g *blockGen) onInputPin(model tin.Model, ptrType r.Type) {
	if g.events != nil {
		// skip registering atoms
		if !g.isMutating() {
			// skip trying to auto-register pins defined by an interface
			// we need a struct to create a block
			if ptrType.Kind() != r.Interface {
				if t, e := model.TypeInfo(ptrType); e != nil {
					g.events.OnError(e)
				} else {
					g.events.OnBlock(t)
				}
			}
		}
	}
}

func (g *blockGen) fieldsOf(containerValue r.Value, containerType r.Type) {
	for i, cnt := 0, containerType.NumField(); i < cnt; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := containerType.Field(i); len(f.PkgPath) == 0 {
			// if the container exists, then we can access the value of the field
			// ( and not just its type )
			var fieldVal r.Value
			if containerValue.IsValid() {
				fieldVal = containerValue.FieldByIndex(f.Index)
			}
			switch f.Name {
			default:
				if item := g.toolboxField(f, fieldVal); item != nil {
					g.block.Items.Append(item)
				}
			case block.NextStatement:
				// process the pin ( the type of the NextStatement field )
				g.onInputPin(tin.MidBlock, f.Type)
				// process the value attached to the pin
				if fieldVal.IsValid() && !fieldVal.IsNil() {
					// struct under the value
					nv := unpackValue(fieldVal)
					// visiting a chain of pointers in a mutation
					if g.isMutating() {
						sub := g.newAtomGenerator()
						sub.fieldsOf(nv, nv.Type())
					} else {
						// visiting a statement block instance
						next := g.genBlock(nv, g.shadowing.Children())
						g.onInputPin(tin.MidBlock, r.PtrTo(nv.Type()))
						g.block.Next = dom.BlockLink{next} // toolbox dom
					}
				}
			}
		}
	}
}

func (g *blockGen) toolboxField(f r.StructField, fieldVal r.Value) (ret dom.Item) {
	itemName := pascal.ToCaps(f.Name)

	// see if the type implements the stringer, ex. an enum.
	if str, ok := asStringer(f.Type, fieldVal); ok {
		ret = g.newField(itemName, str)

	} else {
		switch k := f.Type.Kind(); k {
		default:
			if str, ok := asString(fieldVal); ok {
				ret = g.newField(itemName, str)
			}
		case r.Struct:
			break // for now, just jump out. (ex. block.Option{}s)
		case r.Bool:
			if str, ok := asBool(fieldVal); ok {
				ret = g.newField(itemName, str)
			}

		case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
			if str, ok := asInt(fieldVal); ok {
				ret = g.newField(itemName, str)
			}

		case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
			if str, ok := asUint(fieldVal); ok {
				ret = g.newField(itemName, str)

			}
		case r.Float32, r.Float64:
			if str, ok := asFloat(fieldVal); ok {
				ret = g.newField(itemName, str)
			}

		// input containing another block
		case r.Ptr, r.Interface:
			switch inputType, _ := option.InputOption(f.Tag); inputType {
			case block.ValueInput:
				if b, ok := g.addBlock(tin.TermBlock, fieldVal, f.Type); ok {
					ret = g.newValue(itemName, b)
				}
			case block.StatementInput:
				if b, ok := g.addBlock(tin.MidBlock, fieldVal, f.Type); ok {
					ret = g.newStatement(itemName, b)
				}
			case block.DummyInput:
				if g.isMutating() {
					if g.events != nil {
						e := errutil.New("can't handle mutations inside mutations")
						g.events.OnError(e)
					}
				} else {
					var nv r.Value
					if fieldVal.IsValid() && !fieldVal.IsNil() {
						// write the atom names to a dom.Mutation
						nv = fieldVal.Elem()
					}
					// add the <mutation></mutation> container
					if g.block.Mutation == nil {
						g.block.Mutation = new(dom.BlockMutation)
					}
					// fill with atom names/types for each mutating input
					atoms := newAtoms(g.names, nv, f.Type.Elem())
					if len(atoms) > 0 {
						m := &dom.Mutation{itemName, atoms}
						g.block.Mutation.Append(m)
					}
					// expand the fields directly into the current dom node.
					sub := g.newMutationGenerator(itemName, atoms)
					sub.fieldsOf(nv, f.Type.Elem())
				}
			}
		}
	}
	return
}

// for the toolbox we need unique, predictable input ids for the atoms
// they dont need to be guids here; and they can reset for every block
// or even every mutation if the inputs are scoped to the mutation

func (g *blockGen) addBlock(model tin.Model, fieldVal r.Value, fieldType r.Type) (ret *dom.Block, okay bool) {
	// notify caller we are seeing a new pin
	g.onInputPin(model, fieldType)
	//
	if fieldVal.IsValid() && !fieldVal.IsNil() {
		// notify about the attached value
		nv := unpackValue(fieldVal)
		g.onInputPin(model, r.PtrTo(nv.Type()))
		// visit the attached value
		ret = g.genBlock(nv, g.shadowing.Children())
		okay = true
	}
	return
}

// generate dom.Atom(s) from the passed go values
func newAtoms(names UniqueNames, mval r.Value, mtype r.Type) []*dom.Atom {
	var out []*dom.Atom
	if next, ok := mval, mval.IsValid(); ok {
		for {
			if next, ok = nextField(next); !ok {
				break
			} else {
				atomName := names.GenerateUniqueName()
				atomType := pascal.ToUnderscore(next.Type().Name())
				atom := &dom.Atom{Name: atomName, Type: atomType}
				out = append(out, atom)
			}
		}
	}
	return out
}

// return the value of the passed struct NextStatement
func nextField(structValue r.Value) (ret r.Value, okay bool) {
	next := structValue.FieldByName(block.NextStatement)
	if next.IsValid() && !next.IsNil() {
		ret, okay = unpackValue(next), true
	}
	return
}

// access the underlying type that a pointer or interface references
// usually a struct.
func unpackValue(v r.Value) (ret r.Value) {
	switch k := v.Kind(); k {
	case r.Ptr:
		ret = v.Elem()
	case r.Interface:
		ret = unpackValue(v.Elem())
	default:
		ret = v
	}
	return
}
