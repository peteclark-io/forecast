package ledger

import (
	"bufio"
	"bytes"
)

type lexer struct {
	stream *bufio.Reader
}

func (l *lexer) read() rune {
	ch, _, err := l.stream.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (l *lexer) unread() {
	l.stream.UnreadRune()
}

func (l *lexer) Scan() (Token, string) {
	ch := l.read()

	if ch == eof {
		return EOF, string(ch)
	}

	if isWhitespace(ch) {
		l.unread()
		return WS, l.readType(isWhitespace)
	}

	if isNewLine(ch) {
		return CR, string(ch)
	}

	if isNumber(ch) {
		l.unread()
		return NUMBER, l.readType(isNumber)
	}

	if isDateSeparator(ch) {
		return DATE_SEPARATOR, string(ch)
	}

	if isCleared(ch) {
		return CLEARED_INDICATOR, string(ch)
	}

	if isAccountSeparator(ch) {
		return ACCOUNT_SEPARATOR, string(ch)
	}

	if isLetter(ch) {
		l.unread()
		return IDENT, l.readType(isLetter)
	}

	return ILLEGAL, string(ch)
}

func (l *lexer) readType(check func(ch rune) bool) string {
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	for {
		ch := l.read()

		if ch == eof {
			break
		}

		if !check(ch) {
			l.unread()
			break
		}

		buf.WriteRune(ch)
	}

	return buf.String()
}
