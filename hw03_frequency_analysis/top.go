package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(inputString string) []string {
	topStrings := make([]string, 0)

	if inputString == "" {
		return topStrings
	}

	strInputSlice := strings.Fields(inputString)
	wordCounts := make(map[string]int)

	for _, word := range strInputSlice {
		_, ok := wordCounts[word]
		if !ok {
			wordCounts[word] = 1
		} else {
			wordCounts[word]++
		}
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
	tmpStrings = append(tmpStrings, outPutStr[0].Word)

	for i := 1; i < 10; i++ {
		if outPutStr[i-1].Count != outPutStr[i].Count {
			sort.Strings(tmpStrings)
			topStrings = append(topStrings, tmpStrings...)
			tmpStrings = make([]string, 0)
		}
		tmpStrings = append(tmpStrings, outPutStr[i].Word)

		if i == 9 {
			sort.Strings(tmpStrings)
			topStrings = append(topStrings, tmpStrings...)
		}
	}

	return topStrings
}
