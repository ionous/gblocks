package pascal

import (
	"fmt"
)

func ExampleCamel() {
	fmt.Println(ToCamel("ExampleText"))
	// Output:
	// 	exampleText
}

func ExampleCaps() {
	fmt.Println(ToCaps("ExampleText"))
	// Output:
	// 	EXAMPLE_TEXT
}

func ExampleSpaces() {
	fmt.Println(ToSpaces("ExampleText"))
	// Output:
	// example text
}

func ExampleUnderscore() {
	fmt.Println(ToUnderscore("ExampleText"))
	// Output:
	// 	example_text
}
