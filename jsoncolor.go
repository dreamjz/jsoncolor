package jsoncolor

import (
	"io"
	"os"
	"strconv"

	"github.com/mattn/go-isatty"

	"golang.org/x/term"
)

// Colors specifies colorization of JSON output.
type Colors struct {
	// Null is the color for JSON nil.
	Null Color

	// Bool is the color for boolean values.
	Bool Color

	// Number is the color for number values.
	Number Color

	// String is the color for string values.
	String Color

	// Key is the color for JSON keys.
	Key Color

	// Bytes is the color for byte data.
	Bytes Color

	// Time is the color for datetime values.
	Time Color

	// Punc is the color for JSON punctuation.
	Punc Color
}

// AppendNull appends a colorized "null" to b.
func (c *Colors) AppendNull(b []byte) []byte {
	if c == nil {
		return append(b, "null"...)
	}

	b = append(b, c.Null...)
	b = append(b, "null"...)
	return append(b, ansiReset...)
}

// AppendBool appends the colorized bool v to b.
func (c *Colors) AppendBool(b []byte, v bool) []byte {
	if c == nil {
		if v {
			return append(b, "true"...)
		}

		return append(b, "false"...)
	}

	b = append(b, c.Bool...)
	if v {
		b = append(b, "true"...)
	} else {
		b = append(b, "false"...)
	}

	return append(b, ansiReset...)
}

// AppendKey appends the colorized key v to b.
func (c *Colors) AppendKey(b []byte, v []byte) []byte {
	if c == nil {
		return append(b, v...)
	}

	b = append(b, c.Key...)
	b = append(b, v...)
	return append(b, ansiReset...)
}

// AppendInt64 appends the colorized int64 v to b.
func (c *Colors) AppendInt64(b []byte, v int64) []byte {
	if c == nil {
		return strconv.AppendInt(b, v, 10)
	}

	b = append(b, c.Number...)
	b = strconv.AppendInt(b, v, 10)
	return append(b, ansiReset...)
}

// AppendUint64 appends the colorized uint64 v to b.
func (c *Colors) AppendUint64(b []byte, v uint64) []byte {
	if c == nil {
		return strconv.AppendUint(b, v, 10)
	}

	b = append(b, c.Number...)
	b = strconv.AppendUint(b, v, 10)
	return append(b, ansiReset...)
}

// AppendPunc appends the colorized punctuation mark v to b.
func (c *Colors) AppendPunc(b []byte, v byte) []byte {
	if c == nil {
		return append(b, v)
	}

	b = append(b, c.Punc...)
	b = append(b, v)
	return append(b, ansiReset...)
}

// Color is used to render terminal colors. In effect, Color is
// the ANSI prefix code: the prefix is written, then the actual value,
// then the ANSI reset code.
type Color []byte

// ansiReset is the ANSI ansiReset escape code.
const ansiReset = "\x1b[0m"

// DefaultColors returns the default Colors configuration.
// These colors attempt to follow jq's default colorization.
func DefaultColors() *Colors {
	return &Colors{
		Null:   Color("\x1b[2m"),
		Bool:   Color("\x1b[1m"),
		Number: Color("\x1b[36m"),
		String: Color("\x1b[32m"),
		Key:    Color("\x1b[34;1m"),
		Bytes:  Color("\x1b[2m"),
		Time:   Color("\x1b[32;2m"),
		Punc:   Color("\x1b[1m"),
	}
}

// IsColorTerminal returns true if w is a colorable terminal.
func IsColorTerminal(w io.Writer) bool {
	// This logic could be pretty dodgy; use at your own risk.
	if w == nil {
		return false
	}

	if !isTerminal(w) {
		return false
	}

	if os.Getenv("TERM") == "dumb" {
		return false
	}

	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	if isatty.IsCygwinTerminal(f.Fd()) {
		return false
	}

	return true
}

// isTerminal returns true if w is a terminal.
func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
