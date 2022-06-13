package hw02unpackstring

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	// Place your code here.

	if utf8.ValidString(inputString) {
		return "", ErrInvalidString
	}

	outputString := ""

	for _, rune := range inputString {
		// i может перепрыгивать значения 1,2,4,6,9...
		// r - имеет тип rune, int32

		fmt.Println(rune)
	}

	// склейка строк
	//Join(a []string, sep string) string

	return outputString, nil
}
