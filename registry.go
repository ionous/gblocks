package gblocks

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	r "reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Registry struct {
	types map[string]r.Type
	// given a golang type, find the corresponding blockly field
	// ex. FieldInput	 -> field_input', Blockly.FieldTextInput);
	fields map[r.Type]string
	enums  map[r.Type]enumInfo
}

type stringPair [2]string // display, uniquifier
type enumInfo struct {
	pairs []stringPair
}

var blockPtr = r.TypeOf((*Block)(nil))

func (reg *Registry) RegisterField(name string, f interface{}) {
	if reg.fields == nil {
		reg.fields = make(map[r.Type]string)
	}
	t := r.TypeOf(f).Elem()
	reg.fields[t] = name
}

// New - returns a pointer to passed type name
func (reg *Registry) New(name string) (ret r.Value, err error) {
	if t, ok := reg.types[name]; !ok {
		err = errutil.New("unknown type", name)
	} else {
		ret = r.New(t)
	}
	return
}

func (reg *Registry) RegisterBlock(b interface{}, opt map[string]interface{}) (err error) {
	ptrType := r.TypeOf(b)
	structType := ptrType.Elem()
	name := toTypeName(structType)
	if reg.types == nil {
		reg.types = make(map[string]r.Type)
	} else if _, exists := reg.types[name]; exists {
		panic("type already exists")
	}
	if opt == nil {
		opt = make(map[string]interface{})
	}
	// delay the json initialization, so enums, etc. can all be registered first?
	// useful though to get options out.
	if e := reg.initJson(structType, opt); e != nil {
		err = e
	} else {
		fns := Options{
			"init": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
				b := &Block{Object: obj}

				if e := b.JsonInit(opt); e != nil {
					// hmmm...
				}
				return // init has no return
			}),
		}

		reg.types[name] = structType
		/*
				if m, ok := r.ValueOf(b).Interface().(Mutator); ok {
					// pepare the mutator dialog
					fns["decompose"] = js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
						// for the standard mutator pattern, the first block is either a c-shape or a "head block".
						//  c-shapes contain statement inputs with zero or more elements;
						//  head-blocks start a stack of blocks, each with a prev statement link.
						b := &Block{Object: obj}
						if topBlock, e := b.Workspace.NewBlock(via[0]); e != nil {
							//
						} else {
							topBlock.InitSvg()
							connection := topBlock.NextConnection
							if connection == nil {
								connection = topBlock.GetFirstStatementConnection()
							}
							for i := 0; i < m.NumMutations(); i++ {
								n := m.Mutations(i)
								if item, e := b.Workspace.NewBlock(n.BlockType); e != nil {
									//
								} else {
									item.InitSvg()
									connection.Connect(item.PreviousConnection)
									connection = item.NextConnection
								}
							}
							return topBlock
						}
						return
					})
					// Save the original connections into the mutator dialog blocks.
					// Since the mutator dialog is non-modal, we may have to update the dialog's connections in reponse to changes elsewhere.
					// This allows the mutation's blocks to be reordered, and for this block to be recconstructed accounting for any changes in ordering which may have occured.
					fns["saveConnections"] = js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
						// first block in the mutation dialog
						containerBlock := &Block{Object: obj}
						// we want to traverse *this* and store our connections into the container.
						// we wil have already created the mutation diaog
						// block from the mutation dialog
						// var itemBlock = containerBlock.getInputTargetBlock('STACK');
						// var i = 0;
						// while (itemBlock) {
						//   var input = this.getInput('ADD' + i);
						//   itemBlock.valueConnection_ = input && input.connection.targetConnection;
						//   i++;
						//   itemBlock = itemBlock.nextConnection &&
						//       itemBlock.nextConnection.targetBlock();
						// }
					})

			fns["compose"] = js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
				containerBlock := &Block{Object: obj}
				clauseBlock := containerBlock.NextConnection.TargetBlock()

				for clauseBlock != nil {
					//    switch (clauseBlock.type) {
					//      case 'controls_if_elseif':
					//        this.elseifCount_++;
					//        valueConnections.push(clauseBlock.valueConnection_);
					//        statementConnections.push(clauseBlock.statementConnection_);
					//        break;
					//      case 'controls_if_else':
					//        this.elseCount_++;
					//        elseStatementConnection = clauseBlock.statementConnection_;
					//        break;
					//      default:
					//        throw TypeError('Unknown block type: ' + clauseBlock.type);
					//    }
					//    clauseBlock = clauseBlock.nextConnection &&
					//        clauseBlock.nextConnection.targetBlock();

				//  this.updateShape_();
				//
				m := js.Global.Get("Blockly").Get("Mutator")
				//  // Reconnect any child blocks.
				//  for (var i = 1; i <= this.elseifCount_; i++) {
				// m.Call("reconnect",
				//    Blockly.Mutator.reconnect(valueConnections[i], this, 'IF' + i);
				// m.Call("reconnect",
				//    Blockly.Mutator.reconnect(statementConnections[i], this, 'DO' + i);
				//  }
				// m.Call("reconnect",
				//  Blockly.Mutator.reconnect(elseStatementConnection, this, 'ELSE');
			})
		*/
		//mutationToDom
		//domToMuatation

		// create mutator dialog
		// re/create block from mutator dialog
		//updateShape_
		//

		// the ui system has a mapping of block type names to block prototypes
		// standalone tests do not
		if blocks := js.Global.Get("Blockly").Get("Blocks"); blocks.Bool() {
			blocks.Set(name, fns)
		}
	}
	return
}

// expects a map of intish to string
func (reg *Registry) RegisterEnum(n interface{}) (ret []stringPair, err error) {
	var pairs []stringPair
	if src, srcType := r.ValueOf(n), r.TypeOf(n); srcType.Kind() != r.Map {
		err = errutil.New("invalid enum mapping")
	} else if keyType, valueType := srcType.Key(), srcType.Elem(); valueType.Kind() != r.String {
		err = errutil.New("invalid enum value type")
	} else {
		// want to build an array of display to stringer
		// want to sort that array for display
		// want to store that at the v.Type for lookup.
		// eventually a map? probably of stringer to int ( for reverse conversion, setting in response to changes )
		keys := src.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].Int() < keys[j].Int()
		})

		for _, key := range keys {
			unique := fmt.Sprint(key)
			display := src.MapIndex(key).String()
			pair := stringPair{display, unique}
			pairs = append(pairs, pair)
		}
		if reg.enums == nil {
			reg.enums = make(map[r.Type]enumInfo)
		}
		reg.enums[keyType] = enumInfo{pairs: pairs}
		ret = pairs
	}
	return
}

// note: perhaps anonymous structs could be used to separate into args blocks
func (reg *Registry) initJson(t r.Type, opt Options) (err error) {
	name := toTypeName(t)
	if _, ok := opt[opt_type]; !ok {
		opt[opt_type] = name
	}
	// create "message0": "%1 %2",
	if _, ok := opt[opt_message0]; !ok {
		var els []string
		var arg int
		for i := 0; i < t.NumField(); i++ {
			// skip unexpected symbols ( only unexported symbols have a pkg path )
			if f := t.Field(i); len(f.PkgPath) == 0 {
				// embedded structs represent mutations
				// we dont show mutations by default.
				if f.Type.Kind() != r.Struct && f.Name != previousStatement && f.Name != nextStatement {
					arg += 1 // 1 indexed
					els = append(els, "%"+strconv.Itoa(arg))
				}
			}
		}
		opt[opt_message0] = strings.Join(els, " ")
	}
	//
	if !opt.contains(opt_args0) {
		if args, e := reg.makeArgs(t); e != nil {
			err = errutil.Append(err, e)
		} else if len(args) > 0 {
			opt[opt_args0] = args
		}
	}
	hasPrev, hasNext := opt.contains(opt_previousStatement), opt.contains(opt_nextStatement)
	if !hasPrev {
		if f, ok := t.FieldByName(previousStatement); ok {
			if c, e := constraintsForField(f); e != nil {
				err = errutil.Append(err, e)
			} else {
				opt[opt_previousStatement] = c
				hasPrev = true
			}
		}
	}
	if !hasNext {
		if f, ok := t.FieldByName(nextStatement); ok {
			if c, e := constraintsForField(f); e != nil {
				err = errutil.Append(err, e)
			} else {
				opt[opt_nextStatement] = c
				hasNext = true
			}
		}
	}
	// block can have prev statement or have output
	hasOutput := opt.contains(opt_output)
	if !hasPrev && !hasOutput {
		if m, ok := r.PtrTo(t).MethodByName("Output"); ok {
			if cnt := m.Type.NumOut(); cnt != 1 {
				err = errutil.Append(err, errutil.New("unexpected output count", cnt))
			} else if out := m.Type.Out(0); out.Kind() != r.Interface {
				err = errutil.Append(err, errutil.New("unexpected output type", out.Kind()))
			} else if basicInterface := r.TypeOf((interface{})(nil)); out != basicInterface {
				// basic interface is nil type.
				opt[opt_output] = nil
			} else {
				outName := toTypeName(out)
				if _, ok := reg.types[outName]; !ok {
					err = errutil.Append(err, errutil.New("unknown output type", outName))
				} else {
					opt[opt_output] = outName
				}
			}
		}
	}

	return
}

// return a string, string array, or nil
func constraintsForField(f r.StructField) (ret interface{}, err error) {
	if s, ok := f.Tag.Lookup("check"); ok {
		ret = splitCommas(s)
	} else {
		switch f.Type.Kind() {
		case r.Interface:
			ret = nil
		case r.Ptr:
			ret = toTypeName(f.Type.Elem())
		default:
			err = errutil.New("statement has unexpected type", f.Type)
		}
	}
	return
}

// evaluate the fields of the passed type to generate json usable by blockly initialization
// each field gets its own "options" Options
func (reg *Registry) makeArgs(structType r.Type) (args []Options, err error) {
	for i := 0; i < structType.NumField() && err == nil; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := structType.Field(i); len(f.PkgPath) == 0 &&
			f.Name != previousStatement &&
			f.Name != nextStatement {
			// tags take precedence over type info.
			// ( and some, ex. 'align' come only from tags. )
			opt := parseTags(string(f.Tag))
			// 'name'
			opt.add(opt_name, strings.ToUpper(f.Name))

			// build 'type', 'check'
			switch k := f.Type.Kind(); k {

			// mutation; ignore
			case r.Struct:
				continue

			// array of statements.
			case r.Array:
			case r.Slice:
				switch elType := f.Type.Elem(); elType.Kind() {
				case r.Interface:
					if basicInterface := r.TypeOf((interface{})(nil)); elType != basicInterface {
						opt.add(opt_check, toTypeName(elType))
					}
				case r.Ptr:
					opt.add(opt_check, toTypeName(elType))
				default:
					err = errutil.Append(err, errutil.New(f.Name, "has unexpected type", elType))
				}

			// input containing another block.
			case r.Ptr:
				switch elType := f.Type.Elem(); elType.Kind() {
				case r.Struct:
					opt.add(opt_type, input_value)
					opt.add(opt_check, toTypeName(elType))
				default:
					err = errutil.Append(err, errutil.New(f.Name, "has unexpected type", elType))
				}

			// input of one or more block type
			case r.Interface:
				opt.add(opt_type, input_value)

				if !opt.contains(opt_check) {
					// the basic interface mean no check: it accepts everything
					if basicInterface := r.TypeOf((interface{})(nil)); f.Type != basicInterface {
						var check []string
						// run through all the registered types to see if they implement the interface
						for n, t := range reg.types {
							if t.Implements(f.Type) {
								check = append(check, n)
							}
						}
						if len(check) > 0 {
							opt[opt_check] = check
						}
					}
				}

			// a field of some sort ( ex. angle, checkbox, colour, date, dropdown, image, label, number, text, variable )
			// FIX -- how are we doing variables? always object productions ( ie. a get function )
			default:
				if !opt.contains(opt_type) {
					if fieldType, ok := reg.fields[f.Type]; ok {
						opt[opt_type] = fieldType

					} else if enumType, ok := reg.enums[f.Type]; ok {
						opt[opt_type] = field_dropdown
						opt.add(opt_options, enumType.pairs)

					} else {
						switch f.Type.Kind() {
						case r.Bool:
							opt[opt_type] = field_checkbox

						case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
							opt[opt_type] = field_number //setConstraints
							opt.add(opt_precision, 1)

						case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
							opt[opt_type] = field_number
							opt.add(opt_precision, 1)
							opt.add(opt_min, 0)

						case r.Float32, r.Float64:
							opt[opt_type] = field_number

						case r.String:
							if opt.contains(opt_readOnly) {
								opt[opt_type] = field_label
							} else {
								opt[opt_type] = field_input
							}
							opt.add(opt_text, f.Name)

						default:
							switch r.PtrTo(f.Type) {
							case r.TypeOf((*time.Time)(nil)).Elem():
								opt[opt_type] = field_date
							default:
								err = errutil.New("field has unknown type", f.Name, f.Type)
								break
							}
						}
						// type FieldVariable string
						// type FieldDropdown []string
						// type FieldImageDropdown []FieldImage
						// type FieldImage struct {
						// 	Width, Height int
						// 	Src           string
						// 	Alt           string
						// }
					}
				}
			}
			args = append(args, opt)
		}
	}
	return
}
