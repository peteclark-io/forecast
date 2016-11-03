package ledger_v1

import "unicode"

type Token int

const (
	eof           = rune(0)
	ILLEGAL Token = iota
	EOF
	WS
	CR
	IDENT
	NUMBER
	ACCOUNT_SEPARATOR
	CLEARED_INDICATOR
	DATE_SEPARATOR
	COMMENT
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\v' || ch == '\f'
}

func isNumber(ch rune) bool {
	return unicode.IsNumber(ch) || ch == '.'
}

func isDateSeparator(ch rune) bool {
	return ch == '/'
}

func isCleared(ch rune) bool {
	return ch == '*' || ch == '!'
}

func isAccountSeparator(ch rune) bool {
	return ch == ':'
}

func isNewLine(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '\'')
}
