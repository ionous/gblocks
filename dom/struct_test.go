package dom

import (
	"encoding/xml"
	"fmt"
	"os"
)

func ExampleMarshalIndent() {
	v := &Toolbox{
		Id: "toolbox",
		Categories: []*Category{
			&Category{
				Name: "logic", Colour: "%{BKY_LOGIC_HUE}",
				Categories: []*Category{
					&Category{Name: "If"},
				},
			},
		},
		Blocks: []*Block{
			&Block{
				XMLName: BlockName,
				Type:    "stack_block",
				Next: &Block{
					XMLName: BlockName,
					Type:    "stack_block",
				},
			},
			&Block{
				XMLName: ShadowName,
				Type:    "row_block",
				Items: []Item{
					&Value{
						Name: "INPUT",
						Block: &Block{
							XMLName: BlockName,
							Type:    "row_block",
						},
					},
					&Field{
						Name:    "NUMBER",
						Content: "10",
					},
				},
			},
			&Block{
				XMLName: BlockName,
				Type:    "mutable_block",
				Mutations: &[]*Mutation{
					&Mutation{
						Name: "MUTANT",
						MutableInputs: []*Atom{
							&Atom{"atom_test"},
						},
					},
				},
			},
		},
	}

	output, err := xml.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)
	// Output:
	// 	<xml id="toolbox">
	//  <category name="logic" colour="%{BKY_LOGIC_HUE}">
	//   <category name="If"></category>
	//  </category>
	//  <block type="stack_block">
	//   <next>
	//    <block type="stack_block"></block>
	//   </next>
	//  </block>
	//  <shadow type="row_block">
	//   <value name="INPUT">
	//    <block type="row_block"></block>
	//   </value>
	//   <field name="NUMBER">10</field>
	//  </shadow>
	//  <block type="mutable_block">
	//   <mutation>
	//    <input name="MUTANT">
	//     <atom type="atom_test"></atom>
	//    </input>
	//   </mutation>
	//  </block>
	// </xml>
}
