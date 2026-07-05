package vm

import (
	"fmt"
	"strings"

	"github.com/flily/go-brainfuck/context"
)

const (
	DefaultNewLine = "\n"
)

type ErrorLevel = context.ErrorLevel

const (
	Ignored ErrorLevel = context.Ignored
	Note    ErrorLevel = context.Note
	Remark  ErrorLevel = context.Remark
	Warning ErrorLevel = context.Warning
	Error   ErrorLevel = context.Error
	Fatal   ErrorLevel = context.Fatal
)

type Reason int

const (
	ReasonInvalid Reason = iota
	ReasonHalt
	ReasonCallStackOverflow
	ReasonCallStackEmpty
	ReasonUnsupportedInstruction
	ReasonNoInputDevice
	ReasonNoOutputDevice
	ReasonInputError
	ReasonOutputError
)

var reasonText = map[Reason]string{
	ReasonHalt:                   "halt",
	ReasonCallStackOverflow:      "call stack overflow",
	ReasonCallStackEmpty:         "call stack empty",
	ReasonUnsupportedInstruction: "unsupported instruction",
}

func (r Reason) String() string {
	if text, ok := reasonText[r]; ok {
		return text
	}

	return "unknown"
}

func (r Reason) Error() string {
	return r.String()
}

func (r Reason) RuntimeError(level ErrorLevel, ctx *context.Context, format string, args ...any) *RuntimeError {
	message := fmt.Sprintf(format, args...)
	e := &RuntimeError{
		Reason:  r,
		Level:   level,
		Message: message,
		Note:    "",
		Context: ctx,
	}

	return e
}

func (r Reason) OnNote(ctx *context.Context, format string, args ...any) *RuntimeError {
	return r.RuntimeError(Note, ctx, format, args...)
}

func (r Reason) OnRemark(ctx *context.Context, format string, args ...any) *RuntimeError {
	return r.RuntimeError(Remark, ctx, format, args...)
}

func (r Reason) OnWarning(ctx *context.Context, format string, args ...any) *RuntimeError {
	return r.RuntimeError(Warning, ctx, format, args...)
}

func (r Reason) OnError(ctx *context.Context, format string, args ...any) *RuntimeError {
	return r.RuntimeError(Error, ctx, format, args...)
}

func (r Reason) OnFatal(ctx *context.Context, format string, args ...any) *RuntimeError {
	return r.RuntimeError(Fatal, ctx, format, args...)
}

type RuntimeError struct {
	Reason  Reason
	Level   ErrorLevel
	Message string
	Note    string
	Context *context.Context
}

func (e *RuntimeError) Error() string {
	if e.Context == nil {
		lines := make([]string, 0, 2)

		if len(e.Message) <= 0 {
			lines = append(lines, fmt.Sprintf("[%s]: %s", e.Level, e.Reason))

		} else {
			lines = append(lines, fmt.Sprintf("[%s]: %s", e.Level, e.Message))
		}

		if len(e.Note) > 0 {
			lines = append(lines, fmt.Sprintf("    %s", e.Note))
		}

		return strings.Join(lines, DefaultNewLine)
	}

	message := e.Message
	if len(message) <= 0 {
		message = e.Reason.String()
	}

	title := fmt.Sprintf("%s: %s: %s", e.Context.PositionString(), e.Level, message)
	return title + DefaultNewLine + e.Context.HighlightText("%s", e.Note)
}

func (e *RuntimeError) Is(target error) bool {
	reason, ok := target.(Reason)
	if !ok {
		return false
	}

	return e.Reason == reason
}

func (e *RuntimeError) With(format string, args ...any) *RuntimeError {
	e.Note = fmt.Sprintf(format, args...)
	return e
}
