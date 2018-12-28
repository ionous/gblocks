package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	// r "reflect"
)

// Prepare the mutator's dialog with this block's components.
// just like the original code, this goes from the *data* not the workspace block ui
func (b *Block) decompose(ws *Workspace) (ret *Block, err error) {
	// get the predefined mutator block used for the dialog
	if container, e := ws.NewBlock(b.Type + "$mutation"); e != nil {
		err = e
	} else {
		container.InitSvg()
		// if the block has data already, hook it up to the container
		if ctx := ws.Context(b.Id); ctx.IsValid() {
			srcData := ctx.Elem() // source data
			srcType := srcData.Type()
			// the input names match the field names of the data.
			for i, cnt := 0, container.NumInputs(); i < cnt; i++ {
				in := container.Input(i)
				name := in.Name.FieldName()
				connection := in.Connection
				// connect each data element to the mutation's ui input
				if elsInfo, ok := srcType.FieldByName(name); !ok {
					panic("?")
				} else {
					els := srcData.FieldByIndex(elsInfo.Index)
					//  "cant call Len on Input"
					if cnt := els.Len(); cnt > 0 {
						mutationType := elsInfo.Tag.Get("mutation")
						mutationTypes := ws.reg.mutations[mutationType]
						//
						for i := 0; i < cnt; i++ {
							iface := els.Index(i)
							ptr := iface.Elem()
							el := ptr.Elem()
							//
							elType := toTypeName(el.Type())
							if mutationType, ok := mutationTypes[elType]; !ok {
								e := errutil.New("couldnt find type", elType)
								err = errutil.Append(err, e)
							} else if block, e := ws.NewBlock(mutationType); e != nil {
								err = errutil.Append(err, e)
							} else {
								block.InitSvg()
								connection.Connect(block.PreviousConnection)
								connection = block.NextConnection
							}
						}
					}
				}
			}
		}
		if err == nil {
			ret = container
		} else {
			container.Dispose()
		}
	}
	return
}

// "into" each mutation ui block
// ( found by looking into the container block -- the blocks connected to the inputs )
// store links to the connections of the blocks in the workspace
// ( so that reordering the mutations can re-order the connections )
func (b *Block) saveConnections(ws *Workspace, containerBlock *Block) {
	// for each input in the mutation ui
	for mi, mcount := 0, containerBlock.NumInputs(); mi < mcount; mi++ {
		firstInput := containerBlock.Input(mi)
		if c := firstInput.Connection; c != nil {
			// start with the first block connected to the mutation's ui input
			if itemBlock := c.TargetBlock(); itemBlock != nil {
				// the name of dummy input in the (this) block is the same as the mutation's ui input
				blockInput, baseIndex := b.InputByName(firstInput.Name)
				if m := blockInput.Mutation(); m == nil {
					panic("the input in the block should be a mutation")
				} else {
					// each block connected to the mutation ui's input
					// represents a sub-block in our workspace's block
					// we want to store all of the outgoing connections of that sub-block into mutation's ui block
					itemConnections := itemBlock.ResetConnections()

					for subBlock, subBlocks := 0, m.SubBlocks(); subBlock < subBlocks; subBlock++ {

						for subInput, subInputs := 0, m.SubBlockInputCount(subBlock); subInput < subInputs; subInput++ {
							blockInput := b.Input(baseIndex + subInput)
							itemConnections.StoreConnection(blockInput)
						}
						// next block in the mutation ui
						if c := itemBlock.NextConnection; c != nil {
							itemBlock = c.TargetBlock()
							itemConnections = itemBlock.ResetConnections()
						} else {
							break
						}
					}
				}
			}
		}
	}
}

type ConnectionStore struct {
	*js.Object
	targets *js.Object `js:"targets_"` // array of *Connection
}

func (b *Block) ResetConnections() *ConnectionStore {
	b.store = &ConnectionStore{
		Object:  new(js.Object),
		targets: js.MakeWrapper(make([]*Connection, 0)),
	}
	return b.store
}

func (c *ConnectionStore) StoreConnection(in *Input) {
	var target *Connection
	//blockInput.Connection.TargetConnection
	if in != nil && in.Connection != nil {
		target = in.Connection.TargetConnection
	}
	c.targets.SetIndex(c.targets.Length(), target)
}

// Blockly.Block.prototype.getInputTargetBlock = function(name) {
//   var input = this.getInput(name);
//   return input && input.connection && input.connection.targetBlock();
// };

type ConnectionList struct {
	list []*Connection
}

func (cl *ConnectionList) add(c *Connection) {
	cl.list = append(cl.list, c)
}

func (cl *ConnectionList) contains(c *Connection) (ret bool) {
	for _, oc := range cl.list {
		if oc.Object == c.Object {
			ret = true
		}
	}
	return
}

// re/create the workspace blocks from the mutation dialog ui
// in the original, it counts/saves the connnections ( clause blocks; utemBlocks )
// examning the type to determine the number of inputs
// update shape then creates that number of inputs
// then it walks those inputs in the same order as the saved connnections.
// our update shape builds off of data.
// we want to complelty reset the data so that reconnect events can reconnect the data.

// kinda -- cause in the original the mutator ui basically is creating inputs
// i may have inadveneratly -- in save -- only

// func (b *Block) Compose(ws *Workspace, containerBlock *Block) {
// 	// rebuild the block
// 	ctx := ws.Context(b.Id)
// 	// for each mutation in the mutator ui
// 	for mi, mcount := 0, containerBlock.NumInputs(); mi < mcount; mi++ {
// 		firstInput := containerBlock.Input(mi)
// 		// ugh.
// 		el := ctx.Elem()
// 		elType := el.Type()
// 		elsInfo := elType.FieldByName(firstInput.Name.FieldName())
// 		mutationType := elsInfo.Tag.Get("mutation")
// 		mutationTypes := ws.reg.mutations[mutationType]
// 		//
// 		els := el.FieldByIndex(elsInfo.Index)
// 		out := els.Slice(0, 0)

// 		var connections ConnectionList
// 		if c := firstInput.Connection; c != nil {
// 			// each clause represents one sub-block; one element
// 			for clauseBlock := c.TargetBlock(); clauseBlock != nil; {
// 				// for each clause, create some data
// 				typeName := mutationTypes.findWorkspaceType(clauseBlock.Type)
// 				if v, e := ws.reg.New(typeName); e != nil {
// 					//err = errutil.Append(err, e)
// 					panic(e.Error())
// 				} else {
// 					out = r.Append(out, v)
// 				}

// 				// ex. if the el has 3 inputs; we expect 3 from the associated cause
// 				// for each potential input generated by the el.
// 				// store the connection from the cause block
// 				// could be from the json -- ex. makeArgs
// 				// or, we could even store the results of makeArgs as some info which we generate the json from
// 				// here -- we would get that struct for the workspace type (typeName)
// 				//--- that gives us a number
// 				// does it give us a name???

// 				// ,ight have to dump the block to see what the names look likTest

// 				// options:
// 				// - predict the name of the input
// 				// and store the connection there.s
// 				clauseBlock.store

// 				if c := clauseBlock.NextConnection; c != nil {
// 					clauseBlock = c.TargetBlock()
// 				} else {
// 					break
// 				}
// 			}
// 			els.Set(out)
// 		}
// 	}
// 	// rebuild the workspace ui to create the connectors to re/fill
// 	b.updateShape(ws)
// 	// Reconnect any child blocks.
// 	// for (var i = 0; i < this.itemCount_; i++) {
// 	//   Blockly.Mutator.reconnect(connections[i], this, 'ADD' + i);
// 	// }

// }
