package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/decor"
	"github.com/ionous/gblocks/named"
	r "reflect"
	"strings"
	"time"
)

type Registry struct {
	types     RegisteredTypes
	enums     RegisteredEnums
	mutations RegisteredMutations
	Decor     decor.Decorators
}

func (reg *Registry) RegisterEnum(n interface{}) error {
	_, e := reg.enums.registerEnum(n)
	return e
}

func (reg *Registry) RegisterBlock(b interface{}, blockDesc Dict) error {
	structType := r.TypeOf(b).Elem()
	typeName := named.TypeFromStruct(structType)
	return reg.registerType(typeName, structType, blockDesc)
}

func (reg *Registry) RegisterBlocks(blockDesc map[string]Dict, blocks ...interface{}) (err error) {
	for _, b := range blocks {
		structType := r.TypeOf(b).Elem()
		typeName := named.TypeFromStruct(structType)
		var sub Dict
		if blockDesc, ok := blockDesc[typeName.String()]; ok {
			sub = blockDesc
		}
		if e := reg.registerType(typeName, structType, sub); e != nil {
			e := errutil.New("error registering", typeName, e)
			err = errutil.Append(err, e)
		}
	}
	return
}

func (reg *Registry) RegisterMutation(m interface{}, muiBlocks ...Mutation) (err error) {
	structType := r.TypeOf(m).Elem()
	return reg.mutations.RegisterMutation(structType, muiBlocks...)
}

func (reg *Registry) registerType(typeName named.Type, structType r.Type, blockDesc Dict) (err error) {
	if blockly := GetBlockly(); blockly == nil {
		err = errutil.New("blockly doesnt exist")
	} else {
		if reg.types == nil {
			reg.types = make(RegisteredTypes)
		}
		if !reg.types.RegisterType(structType) {
			err = errutil.New("type already exists " + typeName)
		} else {
			if blockDesc == nil {
				blockDesc = make(Dict)
			}
			if mdesc, e := reg.buildBlockDesc(structType, blockDesc); e != nil {
				err = e
			} else {
				// the ui system has "blocks" -- mapping type names to prototypes; standalone tests do not
				init := js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
					b := &Block{Object: obj}
					// create the blockly block

					if e := b.JsonInit(blockDesc); e != nil {
						panic(e)
					} else if len(mdesc) > 0 {
						// all of the mutation blocks used by all of the mutable inputs in this block
						var quarkNames []named.Type
						for _, md := range mdesc {
							inputName, mutationType := md.input, md.mutation
							if registeredMutation, ok := reg.mutations.GetMutation(mutationType); !ok {
								e := errutil.New("unknown mutation", md)
								panic(e)
							} else if in, index := b.InputByName(inputName); index < 0 {
								e := errutil.New("unknown input", inputName)
								panic(e)
							} else {
								in.ForceMutation(mutationType)
								// add to the palette of types shown by the mutator ui
								quarkNames = append(quarkNames, registeredMutation.quarks...)
							}
						}
						b.SetMutator(NewMutator(quarkNames))
					}
					return
				})
				var fns Dict
				if len(mdesc) == 0 {
					fns = Dict{
						"init": init,
					}
				} else {
					fns = Dict{
						"init": init,
						"mutationToDom": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
							b := &Block{Object: obj}
							dom := b.mutationToDom()
							return dom.Object
						}),
						"domToMutation": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
							b, xmlElement := &Block{Object: obj}, &XmlElement{Object: parms[0]}
							if _, e := b.domToMutation(reg, xmlElement); e != nil {
								panic(e)
							}
							return
						}),
						"decompose": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
							b, mui := &Block{Object: obj}, &Workspace{Object: parms[0]}
							if muiContainer, e := b.decompose(reg, mui); e != nil {
								panic(e)
							} else {
								ret = muiContainer.Object
							}
							return
						}),
						"compose": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
							b, containerBlock := &Block{Object: obj}, &Block{Object: parms[0]}
							if e := b.compose(reg, containerBlock); e != nil {
								panic(e)
							}
							return
						}),
						"saveConnections": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
							b, containerBlock := &Block{Object: obj}, &Block{Object: parms[0]}
							if e := b.saveConnections(containerBlock); e != nil {
								panic(e)
							}
							return
						}),
					}
				}

				// register things to blockly.
				reg.types[typeName] = structType
				blockly.AddBlock(typeName, fns)

				// create custom mutation container
				// the container has one input pin for every mutatable field in the workspace block
				if len(mdesc) > 0 {
					fns := Dict{
						"init": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
							muiContainer := &Block{Object: obj}
							for _, md := range mdesc {
								inputName := md.input
								label := NewFieldLabel(inputName.Friendly(), "")
								in := muiContainer.AppendStatementInput(inputName)
								in.AppendField(label.Field)
								if checks, ok := md.constraints.GetConstraints(); ok {
									in.SetChecks(checks)
								}
							}
							return
						}),
					}
					name := named.SpecialType("mui_container", typeName.String())
					blockly.AddBlock(name, fns)
				}
			}
		}
	}
	return
}

func (reg *Registry) buildCheck(t r.Type, field string, blockDesc Dict, key string) (okay bool, err error) {
	if blockDesc.Contains(key) {
		okay = true
	} else if check, e := reg.types.CheckField(t, field); e != nil {
		err = e
	} else if check, ok := check.GetConstraints(); ok {
		blockDesc[key] = check
		okay = true
	}
	return
}

// could be a map, except maps arent ordered.
type mutationDesc struct {
	input       named.Input //
	mutation    named.Type
	constraints Constraints
}

func (md *mutationDesc) String() string {
	return md.input.String() + ":" + md.mutation.String()
}

func (reg *Registry) buildBlockDesc(t r.Type, blockDesc Dict) (retMui []*mutationDesc, err error) {
	name := named.TypeFromStruct(t)
	blockDesc.Insert(opt_type, name)
	//
	if args, e := reg.buildArgs(t, ""); e != nil {
		err = errutil.Append(err, e)
	} else {
		if msg := args.message(); len(msg) > 0 {
			blockDesc.Insert(opt_message0, msg)
		}
		if len(args.list) > 0 {
			blockDesc.Insert(opt_args0, args.list)
		}
		if !blockDesc.Contains(opt_args0) && !blockDesc.Contains(opt_message0) {
			blockDesc[opt_message0] = name.Friendly()
		}
		retMui = args.mutations
	}
	var hasPrev bool
	if checks, e := reg.buildCheck(t, PreviousField, blockDesc, opt_previous); e != nil {
		err = errutil.Append(err, e)
	} else {
		hasPrev = checks
	}
	if _, e := reg.buildCheck(t, NextField, blockDesc, opt_next); e != nil {
		err = errutil.Append(err, e)
	}
	// block can have prev statement or have output
	hasOutput := blockDesc.Contains(opt_output)
	if !hasPrev && !hasOutput {
		if m, ok := r.PtrTo(t).MethodByName(OutputMethod); ok {
			if cnt := m.Type.NumOut(); cnt != 1 {
				err = errutil.Append(err, errutil.New("unexpected output count", cnt))
			} else {
				out := m.Type.Out(0)
				switch k := out.Kind(); k {
				default:
					err = errutil.Append(err, errutil.New("unexpected output type", k))
				case r.Interface:
					if basicInterface := r.TypeOf((interface{})(nil)); out != basicInterface {
						// basic interface is nil type.
						blockDesc[opt_output] = nil
					} else {
						blockDesc[opt_output] = named.TypeFromStruct(out)
					}
				case r.Ptr:
					blockDesc[opt_output] = named.TypeFromStruct(out.Elem())
				}
			}
		}
	}
	return
}

// evaluate the fields of the passed type to generate json usable by blockly initialization.
// see also: InputMutation.addAtom
func (reg *Registry) buildArgs(t r.Type, path string) (ret *argsOut, err error) {
	var args argsOut
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		if f := t.Field(i); len(f.PkgPath) == 0 && f.Name != PreviousField && f.Name != NextField {
			if e := reg.buildArgDesc(f, path, &args); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	if err == nil {
		ret = &args
	}
	return
}

// return a description suitable for blockly describing the passed field.
// the format doesnt appear to be documented anywhere; derived through examples.
// additional returns an optional mutation input name/pair
func (reg *Registry) buildArgDesc(f r.StructField, path string, out *argsOut) (err error) {
	inputName := named.InputFromField(f)
	typeName := named.TypeFromStruct(f.Type)
	argDesc := parseTags(string(f.Tag))
	itemPath := path + inputName.String()

	// if the field has decoration; add a placeholder label for it.
	if decoration, ok := argDesc[opt_decor].(string); ok {
		// fix? validate dc is a valid decoration
		labelName := named.SpecialType(FieldDecor, decoration, itemPath)
		label := Dict{opt_name: labelName, opt_type: field_label, opt_text: ""}
		out.addArg(label)
	}

	// check for some sort of enumerated type first.
	if enumType, ok := reg.enums.GetEnum(typeName); ok {
		argDesc.Insert(opt_name, itemPath)
		argDesc.Insert(opt_type, field_dropdown)
		argDesc.Insert(opt_options, enumType.pairs)
		out.addArg(argDesc)
	} else {
		// FIX -- how are we doing variables? always object productions ( ie. a get function )
		switch k := f.Type.Kind(); k {

		// slice of statements
		case r.Slice:
			if cs, e := reg.types.CheckType(f.Type.Elem()); e != nil {
				err = errutil.New("invalid slice", inputName, typeName, e)
			} else {
				argDesc.Insert(opt_name, itemPath)
				argDesc.Insert(opt_type, input_statement)
				if cs, ok := cs.GetConstraints(); ok {
					argDesc.Insert(opt_check, cs)
				}
				out.addArg(argDesc)
			}

		// mutation
		case r.Struct:
			mutationType := named.TypeFromStruct(f.Type)
			argDesc.Insert(opt_name, itemPath)
			argDesc.Insert(opt_type, input_dummy)
			argDesc.Insert(opt_mutation, mutationType)
			//
			if md, e := reg.buildMutation(inputName, mutationType, argDesc, out); e != nil {
				err = errutil.New("invalid mutation", inputName, typeName, e)
			} else {
				out.addMutation(md)
			}

		// input containing another block
		case r.Ptr, r.Interface:
			if cs, e := reg.types.CheckType(f.Type); e != nil {
				err = errutil.New("invalid reference", inputName, typeName, e)
			} else {
				argDesc.Insert(opt_name, itemPath)
				argDesc.Insert(opt_type, input_value)
				if cs, ok := cs.GetConstraints(); ok {
					argDesc.Insert(opt_check, cs)
				}
				out.addArg(argDesc)
			}

		// a field of some sort ( ex. angle, checkbox, colour, date, dropdown, image, label, number, text, variable )
		case r.Bool:
			argDesc.Insert(opt_name, itemPath)
			argDesc.Insert(opt_type, field_checkbox)
			out.addArg(argDesc)

		case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
			argDesc.Insert(opt_name, itemPath)
			argDesc.Insert(opt_type, field_number)
			argDesc.Insert(opt_precision, 1)
			out.addArg(argDesc)

		case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
			argDesc.Insert(opt_name, itemPath)
			argDesc.Insert(opt_type, field_number)
			argDesc.Insert(opt_precision, 1)
			argDesc.Insert(opt_min, 0)
			out.addArg(argDesc)

		case r.Float32, r.Float64:
			argDesc.Insert(opt_name, itemPath)
			argDesc.Insert(opt_type, field_number)
			out.addArg(argDesc)

		case r.String:
			var argType string
			if argDesc.Contains(opt_readOnly) {
				argType = field_label
			} else {
				argType = field_input
			}
			argDesc.Insert(opt_name, itemPath)
			argDesc.Insert(opt_type, argType)
			argDesc.Insert(opt_text, f.Name)
			out.addArg(argDesc)

		default:
			switch r.PtrTo(f.Type) {
			case r.TypeOf((*time.Time)(nil)).Elem():
				argDesc.Insert(opt_name, itemPath)
				argDesc.Insert(opt_type, field_date)
				out.addArg(argDesc)

			default:
				err = errutil.New("unknown type", inputName, typeName)
				break
			}
		}
		// type FieldVariable string
		// type FieldImageDropdown []FieldImage
		// type FieldImage struct {
		// 	Width, Height int
		// 	Src           string
		// 	Alt           string
		// }
	}
	return
}

// expand the named type into the args. return a description of mutation.
func (reg *Registry) buildMutation(inputName named.Input, mutationName named.Type, argDesc Dict, out *argsOut) (ret *mutationDesc, err error) {
	if mutation, ok := reg.mutations.GetMutation(mutationName); !ok {
		err = errutil.New(mutation.mutationType, "not registered")
	} else {
		md := mutationDesc{input: inputName, mutation: mutationName}
		for i, cnt := 0, mutation.mutationType.NumField(); i < cnt; i++ {
			switch f := mutation.mutationType.Field(i); {
			// skip unexpected symbols ( only unexported symbols have a pkg path )
			case len(f.PkgPath) > 0:
			case f.Name == PreviousField:
				// error?
			case f.Name == NextField:
				// use the next member to determine the constraints
				if cs, e := reg.types.CheckType(f.Type); e != nil {
					e := errutil.New("invalid reference", e)
					err = errutil.Append(err, e)
				} else {
					// when atoms are added, their inputs are injected after this argument's edge.
					out.addArg(argDesc)
					md.constraints = cs
				}
			default:
				// the fields inside a mutation struct have a path of MUTANT/0/...
				// the first dynamic elements -- which are not generated by the registry -- are MUTANT/1/...
				path := strings.Join([]string{inputName.String(), "0", ""}, "/")
				if e := reg.buildArgDesc(f, path, out); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		if err == nil {
			ret = &md
		}
	}
	return
}
