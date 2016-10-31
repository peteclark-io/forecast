package ledger

type Token int

const (
	eof           = rune(0)
	ILLEGAL Token = iota
	EOF
	WS
	CR
	NUMBER
	IDENT
	DATE_SEPARATOR
	ACCOUNT_SEPARATOR
	CLEARED
	NOT_CLEARED
	COMMENT
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isNumber(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isDateSeparator(ch rune) bool {
	return ch == '/'
}

func isCleared(ch rune) bool {
	return ch == '*'
}

func isNotCleared(ch rune) bool {
	return ch == '!'
}

func isAccountSeparator(ch rune) bool {
	return ch == ':'
}

func isNewLine(ch rune) bool {
	return ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '\'')
}
