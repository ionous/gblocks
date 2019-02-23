package tools

import (
	"encoding/xml"
	"fmt"
	"os"
)

func ExampleMarshalIndent() {
	v := &Toolbox{
		Id: "toolbox",
		Categories: Categories{
			&Category{
				Name: "logic", Colour: "%{BKY_LOGIC_HUE}",
				Categories: Categories{
					&Category{Name: "If"},
				},
			},
		},
		Blocks: Blocks{
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
				Items: Items{
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
				Mutations: &Mutations{
					&Mutation{
						Name: "MUTANT",
						Atoms: Atoms{
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
	//    <atoms name="MUTANT">
	//     <atom type="atom_test"></atom>
	//    </atoms>
	//   </mutation>
	//  </block>
	// </xml>
}
