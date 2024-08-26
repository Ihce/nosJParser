package nosj_parser

import "fmt"

func Deserialize(input string) {
	println("Deserialize")
	fmt.Printf("Last two bytes: %s\n", input[len(input)-2:])
	if input[0:2] != "(<" || input[len(input)-2:] != ")>" {
		panic("There was no root map found within the input.")
	}
}

func parse_map() {
	println("parse_map")

}

func parse_number() {
	println("parse_number")
}

func parse_simple_string() {
	println("parse_simple_string")
}

func parse_complex_string() {
	println("parse_complex_string")
}
