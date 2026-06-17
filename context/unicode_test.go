package context

import (
	"testing"

	"strings"
)

func TestCharWidth(t *testing.T) {
	cases := []struct {
		r   rune
		exp int
	}{
		{r: 0, exp: 0},
		{r: 'a', exp: 1},
		{r: ' ', exp: 1},
		{r: '\t', exp: 8},
		{r: '\v', exp: 1},
		{r: '\f', exp: 1},
		{r: '汉', exp: 2},
		{r: 'あ', exp: 2},
		{r: 'ア', exp: 2},
		{r: 'ｱ', exp: 1},
	}

	for _, c := range cases {
		width := CharWidth(c.r)
		if width != c.exp {
			t.Errorf("expected width %d for rune '%c' (%d), got %d", c.exp, c.r, c.r, width)
			t.Errorf("rune: %c<<<<", c.r)
			t.Errorf("      %s<<<<", strings.Repeat("^", c.exp))
		}
	}
}

func TestStringWidth(t *testing.T) {
	s := "hello\tworld"
	w := StringWidth(s)

	if w != 13 {
		t.Errorf("\n%s<<<<<<\n%s<<<<<<  %d", s, strings.Repeat("^", w), w)
	}
}
