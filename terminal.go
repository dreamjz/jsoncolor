//go:build !windows

package jsoncolor

import (
	"io"
	"os"

	"golang.org/x/term"
)

// IsColorTerminal returns true if w is a colorable terminal.
// It respects [NO_COLOR], [FORCE_COLOR] and TERM=dumb environment variables.
//
// [NO_COLOR]: https://no-color.org/
// [FORCE_COLOR]: https://force-color.org/
func IsColorTerminal(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}

	if w == nil {
		return false
	}

	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	if !term.IsTerminal(int(f.Fd())) {
		return false
	}

	return true
}
