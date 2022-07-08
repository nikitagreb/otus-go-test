package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`[\s";.,!']+`)

const maxOutputLen int = 10

func Top10(inputString string) []string {
	topStrings := make([]string, 0)

	if inputString == "" {
		return topStrings
	}

	strInputSlice := re.Split(inputString, -1)
	wordCounts := make(map[string]int)

	for _, word := range strInputSlice {
		w := strings.ToLower(word)
		if word == "-" {
			continue
		}
		wordCounts[w]++
	}

	outPutStr := make([]struct {
		Word  string
		Count int
	}, 0, len(wordCounts))

	for key, value := range wordCounts {
		outPutStr = append(outPutStr, struct {
			Word  string
			Count int
		}{Word: key, Count: value})
	}

	sort.Slice(outPutStr, func(i, j int) bool {
		return outPutStr[i].Count > outPutStr[j].Count
	})

	tmpStrings := make([]string, 0)
	lenOutPut := len(outPutStr)

	for i := 0; i < lenOutPut; i++ {
		if i != 0 && outPutStr[i-1].Count != outPutStr[i].Count {
			sort.Strings(tmpStrings)
			topStrings = append(topStrings, tmpStrings...)

			if i > maxOutputLen {
				break
			}

			tmpStrings = make([]string, 0)
		}

		tmpStrings = append(tmpStrings, outPutStr[i].Word)

		if i == lenOutPut-1 {
			sort.Strings(tmpStrings)
			topStrings = append(topStrings, tmpStrings...)
		}
	}

	if lenOutPut < maxOutputLen {
		return topStrings[0:lenOutPut]
	}

	return topStrings[0:maxOutputLen]
}
