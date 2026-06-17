package context

import (
	"fmt"
)

type CursorState struct {
	Line   int
	Column int
}

func NewCursorState(line int, column int) *CursorState {
	s := &CursorState{
		Line:   line,
		Column: column,
	}

	return s
}

type Cursor struct {
	CursorState
	File *FileContext
}

func NewCursor(file *FileContext) *Cursor {
	c := &Cursor{
		File: file,
		CursorState: CursorState{
			Line:   0,
			Column: 0,
		},
	}

	return c
}

func NewCursorFromBytes(filename string, data []byte) *Cursor {
	file := ReadFileData(filename, data)
	return NewCursor(file)
}

func NewCursorFromString(filename string, s string) *Cursor {
	return NewCursorFromBytes(filename, []byte(s))
}

func (c *Cursor) Position() (string, int, int) {
	return c.File.Filename, c.Line, c.Column
}

func (c *Cursor) PositionString() string {
	filename, line, column := c.Position()
	return fmt.Sprintf("%s:%d:%d", filename, line+1, column+1)
}

func (c *Cursor) EOFContext() *Context {
	lastLine := c.File.LineContext(c.File.Lines() - 1)
	beginState := NewCursorState(lastLine.Line(), lastLine.Length())
	finishState := NewCursorState(lastLine.Line(), lastLine.Length()+1)
	_, ctx := c.MakeContext(beginState, finishState)
	return ctx
}

func (c *Cursor) CurrentChar() (rune, *Context) {
	current := c.State()
	r, _, _ := c.Rune()
	next := c.PeekState(1)

	_, ctx := c.FinishWith(current, next)
	return r, ctx
}

func (c *Cursor) CurrentLine() (string, *Context) {
	line := c.File.Line(c.Line)
	if line == nil {
		return "", nil
	}

	start := c.State()
	start.Column = 0

	finish := c.State()
	finish.Column = line.Length()

	_, ctx := c.FinishWith(start, finish)
	return line.String(), ctx
}

// Rune returns the rune of current position, and EOL and EOF status
func (c *Cursor) Rune() (rune, bool, bool) {
	return c.Peek(0)
}

// End checks if the cursor is at the end of line and end of file
func (c *Cursor) End() (bool, bool) {
	lastLine := c.Line+1 >= len(c.File.Contents)
	line := c.File.Line(c.Line)

	if c.Column >= len(line.Content) {
		return true, lastLine
	}

	return false, false
}

func (c *Cursor) CurrentEOL() []rune {
	line := c.File.Line(c.Line)
	if line == nil {
		return nil
	}

	return line.EOL
}

func (c *Cursor) PeekState(n int) *CursorState {
	if c.Line >= len(c.File.Contents) {
		return nil
	}

	l := c.File.Line(c.Line)
	if c.Column+n > len(l.Content)+1 {
		return nil
	}

	return NewCursorState(c.Line, c.Column+n)
}

// Peek returns the rune at the offset of current position, and EOL and EOF status.
func (c *Cursor) Peek(n int) (rune, bool, bool) {
	if c.Line >= len(c.File.Contents) {
		return 0, true, true
	}

	l := c.File.Line(c.Line)
	if c.Column+n >= len(l.Content) {
		return 0, true, false
	}

	return l.Content[c.Column+n], false, false
}

func (c *Cursor) PeekString(s string) *Context {
	rs := []rune(s)
	begin := c.State()
	i, r := 0, rune(0)
	for i, r = range rs {
		got, eol, eof := c.Peek(i)
		if got != r || eol || eof {
			return nil
		}
	}

	finish := c.PeekState(i + 1)
	_, ctx := c.FinishWith(begin, finish)
	return ctx
}

func (c *Cursor) next(line int, column int) (int, int, *LineContent) {
	content := c.File.Line(line)
	if content == nil {
		return line, column, nil
	}

	column += 1
	if column > content.Length() {
		c.NextNonEmptyLine()
		column = 0
	}

	return c.Line, column, content
}

func (c *Cursor) nextInLine(line int, column int) (int, int, *LineContent) {
	content := c.File.Line(line)
	if content == nil {
		return line, column, nil
	}

	if c.Column > content.Length() {
		return line, column, nil
	}

	return line, column + 1, content
}

func (c *Cursor) NextInLine() (rune, bool) {
	line, column, content := c.nextInLine(c.Line, c.Column)
	if content == nil {
		return 0, true
	}

	r := content.Rune(column) // checked in nextInLine()
	c.Line = line
	c.Column = column
	return r, column >= content.Length()
}

func (c *Cursor) NextNonEmptyLine() bool {
	line := c.Line + 1
	content := c.File.Line(line)
	for content != nil && content.Length() <= 0 {
		line += 1
		content = c.File.Line(line)
	}

	c.Line = line
	c.Column = 0
	return content == nil
}

func (c *Cursor) Next() (rune, bool, bool) {
	line, column, content := c.next(c.Line, c.Column)
	if content == nil {
		return 0, false, true
	}

	c.Line = line
	c.Column = column
	return c.Rune()
}

func (c *Cursor) NextLine() bool {
	line := c.Line + 1
	content := c.File.Line(line)
	if content == nil {
		return true
	}

	c.Line = line
	c.Column = 0
	return false
}

func (c *Cursor) State() *CursorState {
	return NewCursorState(c.Line, c.Column)
}

func (c *Cursor) SetState(state *CursorState) {
	c.CursorState = *state
}

func (c *Cursor) MakeContext(begin *CursorState, finish *CursorState) (string, *Context) {
	line := c.File.LineContext(begin.Line)
	ctx := line.Mark(begin.Column, finish.Column)
	return ctx.Content(), ctx
}

func (c *Cursor) FinishWith(begin *CursorState, finish *CursorState) (string, *Context) {
	if begin.Line != c.Line || finish.Line != c.Line {
		panic(fmt.Sprintf("cursor context line %d does not match cursor line %d", begin.Line, c.Line))
	}

	return c.MakeContext(begin, finish)
}

func (c *Cursor) FinishTo(offset int) (string, *Context) {
	begin := c.State()
	finish := c.PeekState(offset)
	c.SetState(finish)
	return c.FinishWith(begin, finish)
}

func (c *Cursor) Finish(begin *CursorState) (string, *Context) {
	current := c.State()
	return c.FinishWith(begin, current)
}

func (c *Cursor) NextString(s string) *Context {
	ctx := c.PeekString(s)
	if ctx == nil {
		return nil
	}

	line, column := ctx.Last()
	state := NewCursorState(line, column)
	c.SetState(state)
	return ctx
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

func (c *Cursor) Skip(n int) {
	for i := 0; i < n; i++ {
		_, eof, _ := c.Next()
		if eof {
			return
		}
	}
}

func (c *Cursor) SkipInLine(n int) {
	for i := 0; i < n; i++ {
		_, eol := c.NextInLine()
		if eol {
			return
		}
	}
}

func (c *Cursor) SkipWhitespaceInLine() {
	for {
		r, eol, _ := c.Peek(0)
		if eol || !isWhitespace(r) {
			break
		}

		c.NextInLine()
	}
}

func (c *Cursor) SkipWhitespace() {
	for {
		_, eol, _ := c.Rune()
		if eol {
			if eof := c.NextNonEmptyLine(); eof {
				return
			}

			continue
		}

		r, _, _ := c.Rune()
		if isWhitespace(r) {
			c.Next()

		} else {
			break
		}
	}
}

func (c *Cursor) IsFirstNonWhiteChar() bool {
	line := c.File.Line(c.Line)
	if line == nil {
		return false
	}

	for i := 0; i < c.Column; i++ {
		r := line.Rune(i)
		if !isWhitespace(r) {
			return false
		}
	}

	return true
}
