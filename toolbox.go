package gblocks

import (
	"github.com/ionous/gblocks/named"
	r "reflect"
	"strconv"
	"strings"
)

type Attrs map[string]string

// Toolbox - create a dom element from the passed tag and attrs, and attach the passed content.
// content can include dom nodes or gblocks instance data.
// returns the parent node.
// see also: https://developers.google.com/blockly/guides/configure/web/toolbox
func Toolbox(tag string, attrs map[string]string, content ...interface{}) *XmlElement {
	out := NewXmlElement(tag, attrs)
	return NewTools(out, content...)
}

// NewTools - attach toolbox content to the passed parent.
// see also Toolbox.
func NewTools(parent *XmlElement, content ...interface{}) *XmlElement {
	for _, c := range content {
		if child, ok := c.(*XmlElement); ok {
			parent.AppendChild(child)
		} else {
			kid := NewTool(c)
			parent.AppendChild(kid)
		}
	}
	return parent
}

func NewTool(content interface{}) *XmlElement {
	v := r.ValueOf(content).Elem()
	return toolboxBlock(v, NoShadow)
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
func (s Shadowing) Tag() (ret string) {
	if s == IsShadow {
		ret = "shadow"
	} else {
		ret = "block"
	}
	return
}

// ValueToDom generates xml for the passed gblock instance;
// if useShadowing is true, all child elements will be <shadow> otherwise all elements will be <block>.
func ValueToDom(v r.Value, useShadowing bool) *XmlElement {
	var shadowing Shadowing // default is 0; NoShadowing
	if useShadowing {
		shadowing = SubShadow
	}
	return toolboxBlock(v, shadowing)
}

// returns a <block> (or <shadow>)
func toolboxBlock(v r.Value, shadowing Shadowing) *XmlElement {
	t := v.Type()
	n := named.TypeFromStruct(t)
	el := NewXmlElement(shadowing.Tag(), Attrs{"type": n.String()})
	toolboxFields(el, v, t, shadowing, &toolboxPath{})
	return el
}

// for
type toolboxPath struct {
	parent named.Input // mutation field name
	depth  int
}

func (p *toolboxPath) Next() *toolboxPath {
	return &toolboxPath{p.parent, p.depth + 1}
}

func (p *toolboxPath) IsValid() bool {
	return len(p.parent) > 0
}

func (p *toolboxPath) NewValue(input named.Input) *XmlElement {
	return p.newXmlElement("value", input)
}

func (p *toolboxPath) NewStatement(input named.Input) *XmlElement {
	return p.newXmlElement("statement", input)
}

func (p *toolboxPath) NewField(input named.Input, val string) *XmlElement {
	el := p.newXmlElement("field", input)
	el.SetInnerHTML(val)
	return el
}

func (p *toolboxPath) newXmlElement(tag string, input named.Input) *XmlElement {
	return NewXmlElement(tag, Attrs{"name": p.inputPath(input)})
}

func (p *toolboxPath) inputPath(input named.Input) (ret string) {
	parent, field := p.parent.String(), input.String()
	if len(parent) == 0 {
		ret = field
	} else {
		ret = strings.Join([]string{parent, strconv.Itoa(p.depth), field}, "/")
	}
	return
}

func toolboxFields(el *XmlElement, v r.Value, t r.Type, shadowing Shadowing, path *toolboxPath) {
	var mutationEl *XmlElement
	//
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := t.Field(i); len(f.PkgPath) == 0 {
			switch f.Name {
			case PreviousField:
				// skip
			case NextField:
				// <next>, recursive
				if nv := v.FieldByIndex(f.Index); !nv.IsNil() {
					nv := unpackValue(nv)
					if !path.IsValid() {
						// <next><shadow type='nv.Type.Name()'>...</shadow></next>`
						kid := toolboxBlock(nv, shadowing.Children())
						el.AppendChild(NewXmlElement("next")).AppendChild(kid)
					} else {
						// place the inputs of atoms inside the same block/el parent
						// we want to increase the path name of the atoms' memebers and shadowing depth for sub-blocks.
						toolboxFields(el, nv, nv.Type(), shadowing.Children(), path.Next())
					}
				}
			default:
				input := named.InputFromField(f)
				nv := v.FieldByIndex(f.Index)

				// see if the type implements the stringer, for instance an enum.
				type stringer interface{ String() string }
				if str, ok := nv.Interface().(stringer); ok {
					el.AppendChild(path.NewField(input, str.String()))
				} else {
					switch k := f.Type.Kind(); k {
					case r.Bool:
						field := path.NewField(input, strconv.FormatBool(nv.Bool()))
						el.AppendChild(field)

					case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
						field := path.NewField(input, strconv.FormatInt(nv.Int(), 10))
						el.AppendChild(field)

					case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
						field := path.NewField(input, strconv.FormatUint(nv.Uint(), 10))
						el.AppendChild(field)

					case r.Float32, r.Float64:
						field := path.NewField(input, strconv.FormatFloat(nv.Float(), 'g', -1, 32))
						el.AppendChild(field)

					case r.Struct:
						if path.IsValid() {
							panic("can't handle mutations inside mutations")
						}
						if mutationEl == nil {
							mutationEl = el.AppendChild(NewXmlElement("mutation"))
						}
						// write the atom names
						toolboxMutation(input, nv, mutationEl)
						// expand all of the fields directly into the current node.
						// except for "next" -- which we want to go into value=fieldName at the right spot.
						subPath := &toolboxPath{parent: input}
						toolboxFields(el, nv, nv.Type(), shadowing, subPath)

					// input containing another block
					case r.Ptr, r.Interface:
						if !nv.IsNil() {
							valEl := el.AppendChild(path.NewValue(input))
							kid := toolboxBlock(unpackValue(nv), shadowing.Children())
							valEl.AppendChild(kid)
						}

					case r.Slice:
						if !nv.IsNil() {
							top := el.AppendChild(path.NewStatement(input))
							next := false
							for i, cnt := 0, nv.Len(); i < cnt; i++ {
								if next {
									top = top.AppendChild(NewXmlElement("next"))
								}
								kid := toolboxBlock(unpackValue(nv.Index(i)), shadowing.Children())
								top = top.AppendChild(kid)
								next = true
							}
						}

					default:
						if str := nv.String(); len(str) > 0 {
							field := path.NewField(input, str)
							el.AppendChild(field)
						}
					}
				}
			}
		}
	}
}

func nextField(structValue r.Value) (ret r.Value, okay bool) {
	next := structValue.FieldByName(NextField)
	if next.IsValid() && !next.IsNil() {
		ret, okay = unpackValue(next), true
	}
	return
}

// name is the field name of the mutation struct
func toolboxMutation(name named.Input, mutationStruct r.Value, parent *XmlElement) {
	if next, ok := nextField(mutationStruct); ok {
		atoms := NewXmlElement("atoms", Attrs{"name": name.String()})
		parent.AppendChild(atoms)
		for ; ok; next, ok = nextField(next) {
			typeName := named.TypeFromStruct(next.Type())
			atom := NewXmlElement("atom", Attrs{"type": typeName.String()})
			atoms.AppendChild(atom)
		}
	}
}
