package context

import (
	"github.com/mattn/go-runewidth"
)

var charWidthFixMap = map[rune]int{
	'\t': 8,
	'\v': 1,
	'\f': 1,
}

func CharWidth(r rune) int {
	if width, ok := charWidthFixMap[r]; ok {
		return width
	}

	return runewidth.RuneWidth(r)
}

func CharWidthIn(r rune, index int) int {
	w := CharWidth(r)
	switch r {
	case '\t':
		w = 8 - (index % 8)
	}

	return w
}

func StringWidth(s string) int {
	width := 0
	for _, r := range s {
		width += CharWidthIn(r, width)
	}
	return width
}
