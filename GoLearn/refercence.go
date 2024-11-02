package main

import "fmt"

// 0x1400011e060
// 0x1400011e070
// Output is different because first print is pointer of pointer to name
// Second print is function has input is copy of pointer and show new pointer of pointer
func main() {
	name := "bill"

	namePointer := &name

	fmt.Println(&namePointer)
	printPointer(namePointer)
}

func printPointer(namePointer *string) {
	//fmt.Println(namePointer)
	fmt.Println(&namePointer)
}
