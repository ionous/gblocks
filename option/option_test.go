package option

import (
	"fmt"
)

func ExampleMessage() {
	fmt.Println(Message(0), Message(1))
	// Output:
	// 	message0 message1
}

func ExampleArgs() {
	fmt.Println(Args(0), Args(1))
	// Output:
	// 	args0 args1
}
