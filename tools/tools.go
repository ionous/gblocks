package tools

import (
	"encoding/xml"
	"github.com/ionous/gblocks/block"
	r "reflect"
	"strconv"
	"strings"
)

// Toolbox - create a dom element from the passed tag and attrs, and attach the passed content.
// content can include dom nodes or gblocks instance data.
// returns the parent node.
// see also: https://developers.google.com/blockly/guides/configure/web/toolbox
func NewToolbox(contents ...interface{}) *Toolbox {
	var tb Toolbox
	tb.AddBlocks(contents...)
	return &tb
}

// AddBlocks - attach toolbox content to the passed parent.
// see also Toolbox.
func (tb *Toolbox) AddBlocks(contents ...interface{}) {
	for _, c := range contents {
		v := r.ValueOf(c).Elem()
		b := toolboxBlock(v, NoShadow)
		tb.Blocks = append(tb.Blocks, b)
	}
}

// access the underlying value that a pointer or interface refers to
func unpackType(t r.Type) (ret r.Type) {
	switch k := t.Kind(); k {
	case r.Ptr:
		ret = t.Elem()
	case r.Interface:
		ret = unpackType(t.Elem())
	default:
		ret = t
	}
	return
}

// access the underlying type that a pointer or interface references
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

// Shadowing - when creating xml from golang types should we create shadow blocks
// https://developers.google.com/blockly/guides/configure/web/toolbox#shadow_blocks
type Shadowing int

const (
	NoShadow Shadowing = iota
	IsShadow
	SubShadow
)

// Children of shadows or subshadows are shadows
func (s Shadowing) Children() Shadowing {
	if s == SubShadow {
		s = IsShadow // upgrade shadowing; otherwise no change
	}
	return s
}

// Tag is either <shadow> or <block> depending
func (s Shadowing) Tag() (ret xml.Name) {
	if s == IsShadow {
		ret = ShadowName
	} else {
		ret = BlockName
	}
	return
}

// returns a <block> (or <shadow>)
func toolboxBlock(v r.Value, shadowing Shadowing) *Block {
	t := v.Type()
	n := block.TypeFromStruct(t)
	el := &Block{XMLName: shadowing.Tag(), Type: n.String()}
	toolboxFields(el, v, t, shadowing, &toolboxPath{})
	return el
}

//
type toolboxPath struct {
	parent block.Item // mutation field name
	depth  int
}

func (p *toolboxPath) Next() *toolboxPath {
	return &toolboxPath{p.parent, p.depth + 1}
}

func (p *toolboxPath) IsValid() bool {
	return len(p.parent) > 0
}

func (p *toolboxPath) NewValue(input block.Item) *Value {
	return &Value{Name: p.inputPath(input)}
}

func (p *toolboxPath) NewStatement(input block.Item) *Statement {
	return &Statement{Name: p.inputPath(input)}
}

func (p *toolboxPath) NewField(input block.Item, val string) *Field {
	field := &Field{Name: p.inputPath(input)}
	field.Content = val
	return field
}

func (p *toolboxPath) inputPath(input block.Item) (ret string) {
	parent, field := p.parent.String(), input.String()
	if len(parent) == 0 {
		ret = field
	} else {
		ret = strings.Join([]string{parent, strconv.Itoa(p.depth), field}, "/")
	}
	return
}

func toolboxFields(el *Block, v r.Value, t r.Type, shadowing Shadowing, path *toolboxPath) {
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := t.Field(i); len(f.PkgPath) == 0 {
			switch f.Name {
			case block.PreviousStatement:
				// skip
			case block.NextStatement:
				// <next>, recursive
				if nv := v.FieldByIndex(f.Index); !nv.IsNil() {
					nv := unpackValue(nv)
					if !path.IsValid() {
						// <next><shadow type='nv.Type.Name()'>...</shadow></next>`
						el.Next = toolboxBlock(nv, shadowing.Children())
					} else {
						// place the inputs of atoms inside the same block/el parent
						// we want to increase the path name of the atoms' memebers and shadowing depth for sub-blocks.
						toolboxFields(el, nv, nv.Type(), shadowing.Children(), path.Next())
					}
				}
			default:
				itemName := block.ItemFromField(f)
				nv := v.FieldByIndex(f.Index)

				// see if the type implements the stringer, for instance an enum.
				type stringer interface{ String() string }
				if str, ok := nv.Interface().(stringer); ok {
					el.AppendItem(path.NewField(itemName, str.String()))
				} else {
					switch k := f.Type.Kind(); k {
					case r.Bool:
						field := path.NewField(itemName, strconv.FormatBool(nv.Bool()))
						el.AppendItem(field)

					case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
						field := path.NewField(itemName, strconv.FormatInt(nv.Int(), 10))
						el.AppendItem(field)

					case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
						field := path.NewField(itemName, strconv.FormatUint(nv.Uint(), 10))
						el.AppendItem(field)

					case r.Float32, r.Float64:
						field := path.NewField(itemName, strconv.FormatFloat(nv.Float(), 'g', -1, 32))
						el.AppendItem(field)

					case r.Struct:
						if path.IsValid() {
							panic("can't handle mutations inside mutations")
						}
						if el.Mutations == nil {
							el.Mutations = new(Mutations)
						}

						// write the atom names
						toolboxMutation(el, itemName, nv)
						// expand all of the fields directly into the current node.
						// except for "next" -- which we want to go into value=fieldName at the right spot.
						subPath := &toolboxPath{parent: itemName}
						toolboxFields(el, nv, nv.Type(), shadowing, subPath)

					// itemName containing another block
					case r.Ptr, r.Interface:
						if !nv.IsNil() {
							valEl := path.NewValue(itemName)
							valEl.Block = toolboxBlock(unpackValue(nv), shadowing.Children())
							el.AppendItem(valEl)
						}

					case r.Slice:
						if !nv.IsNil() {
							var next *Block
							top := path.NewStatement(itemName)
							for i, cnt := 0, nv.Len(); i < cnt; i++ {
								kid := toolboxBlock(unpackValue(nv.Index(i)), shadowing.Children())
								if next != nil {
									next.Next, next = kid, kid
								} else {
									top.Block, next = kid, kid
								}
							}
							el.AppendItem(top)
						}

					default:
						if str := nv.String(); len(str) > 0 {
							field := path.NewField(itemName, str)
							el.AppendItem(field)
						}
					}
				}
			}
		}
	}
}

// return the value of the passed struct's .NextStatement
func nextField(structValue r.Value) (ret r.Value, okay bool) {
	next := structValue.FieldByName(block.NextStatement)
	if next.IsValid() && !next.IsNil() {
		ret, okay = unpackValue(next), true
	}
	return
}

// name is the field name of the mutation struct
func toolboxMutation(el *Block, name block.Item, mutationStruct r.Value) {
	if next, ok := nextField(mutationStruct); ok {
		atoms := &Mutation{Name: name.String()}
		for ; ok; next, ok = nextField(next) {
			typeName := block.TypeFromStruct(next.Type())
			atoms.AppendAtom(typeName.String())
		}
		el.AppendMutation(atoms)
	}
}
