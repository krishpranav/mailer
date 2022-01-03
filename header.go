// headers.go implements "Q" encoding as specified by RFC 2047.
//Modified from https://github.com/joegrasse/mime
package mailer

import (
	"bufio"
	"io"
	"strings"
)

type encoder struct {
	w         *bufio.Writer
	charset   string
	usedChars int
}

func NewEncoder(w io.Writer, c string, u int) *encoder {
	return &encoder{bufio.NewWriter(w), strings.ToUpper(c), u}
}

func isVchar(c byte) bool {
	return '!' <= c && c <= '~'
}

func isWSP(c byte) bool {
	return c == ' ' || c == '\t'
}
