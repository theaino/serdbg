package main

import (
	"math"
	"strings"
)

func WrapString(text string, width int) (string, int) {
	lineCount := int(math.Ceil(float64(len(text)) / float64(width)))
	lines := make([]string, lineCount)
	for idx := range lineCount {
		start := idx * width
		end := start + width
		lines[idx] = text[start:min(end, len(text))]
	}
	return strings.Join(lines, "\n"), lineCount
}

func JoinStringWrapped(elems []string, sep string, width int) (string, int) {
	lines := make([]string, 1)
	for _, elem := range elems {
		lastLine := lines[len(lines)-1]
		joinedLine := elem
		if lastLine != "" {
			joinedLine = lastLine + sep + joinedLine
		}
		if len(joinedLine) > width {
			lines = append(lines, elem)
		} else {
			lines[len(lines)-1] = joinedLine
		}
	}
	return strings.Join(lines, "\n"), len(lines)
}
