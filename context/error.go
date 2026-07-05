package context

import (
	"fmt"
	"strings"
)

type ErrorLevel int

const (
	Ignored ErrorLevel = iota
	Note
	Remark
	Warning
	Error
	Fatal
)

var levelNames = map[ErrorLevel]string{
	Ignored: "ignored",
	Note:    "note",
	Remark:  "remark",
	Warning: "warning",
	Error:   "error",
	Fatal:   "fatal",
}

func (l ErrorLevel) String() string {
	return levelNames[l]
}

func (l ErrorLevel) NewDiagnostic(ctx *Context, message string, note string) *Diagnostic {
	return NewDiagnostic(l, ctx, message, note)
}

type DiagnosticInfo interface {
	error
	Level() ErrorLevel
}

/*
 * filename:line:column: level: message
 *  line | code
 *       | ^^^^
 *       | note
 */
type Diagnostic struct {
	level   ErrorLevel
	context *Context
	message string
	note    string
}

func NewDiagnostic(level ErrorLevel, ctx *Context, message string, note string) *Diagnostic {
	e := &Diagnostic{
		level:   level,
		context: ctx,
		message: message,
		note:    note,
	}

	return e
}

func NewNote(ctx *Context, message string, args ...any) *Diagnostic {
	m := fmt.Sprintf(message, args...)
	return Note.NewDiagnostic(ctx, m, "")
}

func NewWarning(ctx *Context, message string, args ...any) *Diagnostic {
	m := fmt.Sprintf(message, args...)
	return Warning.NewDiagnostic(ctx, m, "")
}

func NewError(ctx *Context, message string, args ...any) *Diagnostic {
	m := fmt.Sprintf(message, args...)
	return Error.NewDiagnostic(ctx, m, "")
}

func (d *Diagnostic) Level() ErrorLevel {
	return d.level
}

func (d *Diagnostic) Error() string {
	messageLine := fmt.Sprintf("%s: %s: %s", d.context.PositionString(), d.level, d.message)
	return messageLine + DefaultNewLine + d.context.HighlightText("%s", d.note)
}

func (d *Diagnostic) With(note string, args ...any) *Diagnostic {
	d.note = fmt.Sprintf(note, args...)
	return d
}

func (d *Diagnostic) For(reason *Diagnostic) *DiagnosticCombo {
	c := &DiagnosticCombo{
		info:   d,
		reason: reason,
	}

	return c
}

func (d *Diagnostic) ToContainer() *DiagnosticContainer {
	c := NewDiagnosticContainer(d.level)
	_ = c.Add(d)
	return c
}

type DiagnosticCombo struct {
	info   *Diagnostic
	reason *Diagnostic
}

func (c *DiagnosticCombo) Level() ErrorLevel {
	return c.info.level
}

func (c *DiagnosticCombo) Error() string {
	parts := make([]string, 0, 2)
	parts = append(parts, c.info.Error())

	if c.reason != nil {
		parts = append(parts, c.reason.Error())
	}

	return strings.Join(parts, DefaultNewLine)
}

func (c *DiagnosticCombo) ToContainer() *DiagnosticContainer {
	container := NewDiagnosticContainer(c.info.level)
	_ = container.Add(c)
	return container
}

type DiagnosticContainer struct {
	Diagnostics []DiagnosticInfo
	RaiseLevel  ErrorLevel
}

func NewDiagnosticContainer(level ErrorLevel) *DiagnosticContainer {
	c := &DiagnosticContainer{
		Diagnostics: make([]DiagnosticInfo, 0, 8),
		RaiseLevel:  level,
	}

	return c
}

func (c *DiagnosticContainer) Add(d DiagnosticInfo) error {
	c.Diagnostics = append(c.Diagnostics, d)

	if d.Level() >= c.RaiseLevel {
		return d
	}

	return nil
}

func (c *DiagnosticContainer) Merge(other *DiagnosticContainer) error {
	var err error
	for _, d := range other.Diagnostics {
		if e := c.Add(d); e != nil && err == nil {
			err = e
		}
	}

	return err
}

func (c *DiagnosticContainer) Level() ErrorLevel {
	level := Ignored
	for _, d := range c.Diagnostics {
		if d.Level() > level {
			level = d.Level()
		}
	}

	return level
}

func (c *DiagnosticContainer) Error() string {
	parts := make([]string, 0, len(c.Diagnostics))
	for _, d := range c.Diagnostics {
		parts = append(parts, d.Error())
	}

	return strings.Join(parts, DefaultNewLine)
}

func (c *DiagnosticContainer) Count(level ErrorLevel) int {
	count := 0
	for _, d := range c.Diagnostics {
		if d.Level() >= level {
			count++
		}
	}

	return count
}
