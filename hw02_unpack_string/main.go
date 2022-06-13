package hw02unpackstring

import (
	"fmt"
)

func main() {

	var inputString string
	fmt.Println("Please enter a line")
	_, err := fmt.Scanf("%s\n", &inputString)
	if err != nil {
		fmt.Println("error scan string")
	}

	_, _ = Unpack(inputString)
}
