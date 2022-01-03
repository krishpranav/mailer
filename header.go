// headers.go implements "Q" encoding as specified by RFC 2047.
//Modified from https://github.com/joegrasse/mime to use with Go Simple Mail

package mailer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type encoder struct {
	w         *bufio.Writer
	charset   string
	usedChars int
}

func newEncoder(w io.Writer, c string, u int) *encoder {
	return &encoder{bufio.NewWriter(w), strings.ToUpper(c), u}
}

func (e *encoder) encode(p []byte) (n int, err error) {
	var output bytes.Buffer
	allPrintable := true

	maxLineLength := 76

	p = secureHeader(p)

	for _, c := range p {
		if !isVchar(c) && !isWSP(c) {
			allPrintable = false
			break
		}
	}

	if allPrintable {
		text := string(p)
		words := strings.Split(text, " ")

		lineBuffer := ""
		firstWord := true

		for _, word := range words {
			newWord := ""
			if !firstWord {
				newWord += " "
			}
			newWord += word

			if (e.usedChars+len(lineBuffer)+len(newWord)) > maxLineLength && (lineBuffer != "" || e.usedChars != 0) {
				output.WriteString(lineBuffer + "\r\n")
				if !firstWord {
					lineBuffer = ""
				} else {
					lineBuffer = " "
				}

				e.usedChars = 0
			}

			lineBuffer += newWord

			firstWord = false

		}

		output.WriteString(lineBuffer)

	} else {
		firstLine := true
		if e.usedChars == 0 {
			maxLineLength = 75
		}

		wordBegin := "=?" + e.charset + "?Q?"
		wordEnd := "?="

		lineBuffer := wordBegin

		for i := 0; i < len(p); {
			encodedChar, runeLength := encode(p, i)

			if len(lineBuffer)+e.usedChars+len(encodedChar) > (maxLineLength - len(wordEnd)) {
				output.WriteString(lineBuffer + wordEnd + "\r\n")
				lineBuffer = " " + wordBegin
				firstLine = false
			}

			lineBuffer += encodedChar

			i = i + runeLength

			if !firstLine {
				e.usedChars = 0
				maxLineLength = 76
			}
		}

		output.WriteString(lineBuffer + wordEnd)
	}

	e.w.Write(output.Bytes())
	e.w.Flush()
	n = output.Len()

	return n, nil
}

func encode(text []byte, i int) (encodedString string, runeLength int) {
	started := false

	for ; i < len(text) && (!utf8.RuneStart(text[i]) || !started); i++ {
		switch c := text[i]; {
		case c == ' ':
			encodedString += "_"
		case isVchar(c) && c != '=' && c != '?' && c != '_':
			encodedString += string(c)
		default:
			encodedString += fmt.Sprintf("=%02X", c)
		}

		runeLength++

		started = true
	}

	return
}

func secureHeader(text []byte) []byte {
	secureValue := strings.TrimSpace(string(text))
	secureValue = strings.Replace(secureValue, "\r", "", -1)
	secureValue = strings.Replace(secureValue, "\n", "", -1)
	secureValue = strings.Replace(secureValue, "\t", "", -1)

	return []byte(secureValue)
}

func isVchar(c byte) bool {
	return '!' <= c && c <= '~'
}

func isWSP(c byte) bool {
	return c == ' ' || c == '\t'
}
