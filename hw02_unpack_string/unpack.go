package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type directoryData struct {
	Numbers       map[rune]int8
	Slash         rune
	SlashNumber   int32
	LetterNNumber int32
	SlashQuote    rune
}

func Unpack(inputString string) (string, error) {
	directory := initDirectoryData()

	if strings.Contains(inputString, "10") || checkFirstCharNumber(inputString, directory) {
		return "", ErrInvalidString
	}

	if inputString == "" {
		return "", nil
	}

	prevSlash := false
	prevSecondSlash := false
	prevThirdSlash := false
	prevLineBreak := false
	var checkRune bool
	var letters []string

	for _, runeValue := range inputString {
		checkRune = false
		if isNumber(runeValue, directory) && prevLineBreak && !checkRune {
			letters = copySlashN(runeValue, directory, letters)
			checkRune = true
		}

		if runeValue != directory.SlashNumber && prevSlash && prevSecondSlash && prevThirdSlash && !checkRune {
			letters = letters[:len(letters)-1]
			letters = letters[:len(letters)-1]
			letters = append(letters, string(runeValue))
			checkRune = true
		}

		if runeValue == '0' && !checkRune {
			letters = letters[:len(letters)-1]
			checkRune = true
		}

		if runeValue == '1' {
			continue
		}

		if isNumber(runeValue, directory) && prevSlash && prevSecondSlash && !checkRune {
			letters = letters[:len(letters)-1]
			letters = copyRune(runeValue, directory, letters)
			checkRune = true
		}

		if isNumber(runeValue, directory) && prevSlash && !checkRune {
			letters = letters[:len(letters)-1]
			letters = append(letters, string(runeValue))
			checkRune = true
		}

		if isNumber(runeValue, directory) && !checkRune {
			letters = copyRune(runeValue, directory, letters)
			checkRune = true
		}

		if !checkRune {
			letters = append(letters, string(runeValue))
		}

		prevThirdSlash = prevSecondSlash
		prevSecondSlash = prevSlash
		prevLineBreak = runeValue == directory.LetterNNumber && prevSlash
		prevSlash = runeValue == directory.SlashNumber
	}

	return strings.Join(letters, ""), nil
}

func initDirectoryData() directoryData {
	return directoryData{
		Numbers: map[rune]int8{
			'0': 0,
			'1': 1,
			'2': 2,
			'3': 3,
			'4': 4,
			'5': 5,
			'6': 6,
			'7': 7,
			'8': 8,
			'9': 9,
		},
		Slash:         '\\',
		SlashQuote:    '`',
		SlashNumber:   92,
		LetterNNumber: 110,
	}
}

func isNumber(sRune rune, directory directoryData) bool {
	for iRune := range directory.Numbers {
		if sRune == iRune {
			return true
		}
	}

	return false
}

func copySlashN(runeValue rune, directory directoryData, letters []string) []string {
	var i int8 = 2
	for ; i <= directory.Numbers[runeValue]; i++ {
		letters = append(letters, string(directory.SlashNumber))
		letters = append(letters, string(directory.LetterNNumber))
	}

	return letters
}

func copyRune(runeValue rune, directory directoryData, letters []string) []string {
	var i int8 = 2
	for ; i <= directory.Numbers[runeValue]; i++ {
		letters = append(letters, letters[len(letters)-1])
	}

	return letters
}

func checkFirstCharNumber(inputString string, directory directoryData) bool {
	for runeIndex, runeValue := range inputString {
		if runeIndex == 0 && isNumber(runeValue, directory) {
			return true
		}
	}

	return false
}
