package envy_test

import envy "."
import "fmt"
import "time"

/*
This may not pass if go decides to execute
the slots in different order
*/
func ExampleMain() {
	envy.Connect(ExampleSignal, ExampleSlot1)
	envy.Connect(ExampleSignal, ExampleSlot2)
	envy.Emit(ExampleSignal, 42, "cool things")
	time.Sleep(1*time.Second)
	// Output:
	// 42 is a cool number!
	// I have 42 cool things!
}

var ExampleSignal envy.Signal = envy.New()

func ExampleSlot1(n int) {
	fmt.Printf("%d is a cool number!\n", n)
}

func ExampleSlot2(n int, thing string) {
	fmt.Printf("I have %d %s!\n", n, thing)
}