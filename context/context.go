package context

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	DefaultNewLine = "\n"
)

type ByLineContextLine []*LineContext

func (a ByLineContextLine) Len() int {
	return len(a)
}

func (a ByLineContextLine) Less(i, j int) bool {
	return a[i].Content.Line < a[j].Content.Line
}

func (a ByLineContextLine) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type Context struct {
	File      *FileContext
	PrevLines []*LineContext
	NextLines []*LineContext
	Lines     []*LineContext
}

type ContextProvider interface {
	Context() *Context
}

func (c *Context) Join(ctxs ...*Context) *Context {
	lines := make([]*LineContext, 0, len(c.Lines)+len(ctxs)*10)
	for _, line := range c.Lines {
		lines = append(lines, line.Duplicate())
	}

	for i, ctx := range ctxs {
		if ctx == nil {
			continue
		}

		if ctx.File != c.File {
			err := fmt.Errorf("context %d does not match file", i)
			panic(err)
		}

		for _, l := range ctx.Lines {
			line := FindLineContextSameLine(lines, l)
			if line == nil {
				lines = append(lines, l.Duplicate())

			} else {
				line.MergeHighlights(l)
			}
		}
	}

	result := &Context{
		File:  c.File,
		Lines: make([]*LineContext, 0, len(lines)),
	}

	for _, line := range lines {
		sort.Sort(ByHighlight(line.Highlights))
		result.Lines = append(result.Lines, line)
	}

	sort.Sort(ByLineContextLine(result.Lines))
	return result
}

func (c *Context) Load(prev int, next int) {
	first := true
	lineFirst, lineLast := 100000000000, 0
	for _, line := range c.Lines {
		if first {
			lineFirst, lineLast = line.Content.Line, line.Content.Line

		} else {
			if line.Content.Line < lineFirst {
				lineFirst = line.Content.Line
			}

			if line.Content.Line > lineLast {
				lineLast = line.Content.Line
			}
		}

		first = false
	}

	c.PrevLines = make([]*LineContext, 0, prev)
	c.NextLines = make([]*LineContext, 0, next)
	for i := 0; i < len(c.File.Contents); i++ {
		np, nl := lineFirst-i, i-lineLast
		if np > 0 && np <= prev {
			c.PrevLines = append(c.PrevLines, c.File.LineContext(i))
		}

		if nl > 0 && nl <= next {
			c.NextLines = append(c.NextLines, c.File.LineContext(i))
		}
	}
}

func (c *Context) HighlightTextWith(indicator string, format string, args ...any) string {
	parts := make([]string, 0, len(c.Lines)+len(c.PrevLines)+len(c.NextLines)+3)

	for _, line := range c.PrevLines {
		parts = append(parts, line.String())
	}

	for i, line := range c.Lines {
		if i == len(c.Lines)-1 {
			// show message in last line
			parts = append(parts, line.HighlighTextWith(indicator, format, args...))
		} else {
			parts = append(parts, line.HighlighText(NoHighlightMessage))
		}
	}

	for _, line := range c.NextLines {
		parts = append(parts, line.String())
	}

	return strings.Join(parts, DefaultNewLine)
}

func (c *Context) HighlightText(format string, args ...any) string {
	return c.HighlightTextWith(DefaultIndicator, format, args...)
}

func (c *Context) Last() (int, int) {
	line, column := 0, 0
	for _, lineCtx := range c.Lines {
		lineNo := lineCtx.Content.Line
		if lineNo >= line {
			line, column = lineNo, lineCtx.Last()
		}
	}

	return line, column
}

func (c *Context) Content() string {
	if len(c.Lines) <= 0 {
		return ""
	}

	return c.Lines[0].HighlightContent()
}

func (c *Context) NextInLineContext() *Context {
	line, column := c.Last()
	lineCtx := c.File.LineContext(line)
	if lineCtx == nil {
		return nil
	}

	if column < lineCtx.Length() {
		result := lineCtx.Mark(column, column+1)
		return result
	}

	result := lineCtx.Mark(column, column+1)
	return result
}

func (c *Context) NextContext() *Context {
	line, column := c.Last()
	lineCtx := c.File.LineContext(line)
	if lineCtx == nil {
		return nil
	}

	if column < lineCtx.Length() {
		result := lineCtx.Mark(column, column+1)
		return result
	}

	nextLineCtx := c.File.LineContext(line + 1)
	if nextLineCtx == nil {
		result := lineCtx.Mark(column, column)
		return result
	}

	result := nextLineCtx.Mark(0, 1)
	return result
}

func (c *Context) Note(format string, args ...any) *Diagnostic {
	return NewNote(c, format, args...)
}

func (c *Context) Warning(format string, args ...any) *Diagnostic {
	return NewWarning(c, format, args...)
}

func (c *Context) Error(format string, args ...any) *Diagnostic {
	return NewError(c, format, args...)
}

func (c *Context) Position() (string, int, int) {
	filename := c.File.Filename
	line, column := c.Lines[0].Position()
	return filename, line, column
}

func (c *Context) PositionString() string {
	filename, line, column := c.Position()
	return fmt.Sprintf("%s:%d:%d", filename, line+1, column+1)
}

func Join(ctxs ...*Context) *Context {
	if len(ctxs) == 0 {
		return nil
	}

	firstNonNil := 0
	for i := 0; i < len(ctxs); i++ {
		if ctxs[i] != nil {
			firstNonNil = i
			break
		}
	}

	first := ctxs[firstNonNil]
	return first.Join(ctxs[firstNonNil+1:]...)
}

func JoinObjects(objs ...ContextProvider) *Context {
	ctxs := make([]*Context, 0, len(objs))
	for _, obj := range objs {
		val := reflect.ValueOf(obj)
		if val.Kind() == reflect.Invalid || (val.Kind() == reflect.Pointer && val.IsNil()) {
			continue
		}
		ctxs = append(ctxs, obj.Context())
	}

	return Join(ctxs...)
}
