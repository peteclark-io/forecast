package ledger

import (
	"bytes"
	"unicode"
)

type Token int

const (
	eof           = rune(0)
	ILLEGAL Token = iota
	SOF
	EOF
	SPACE
	BLOCK_SPACE
	CR_BLOCK_SPACE
	CR
	ACCOUNT
	PAYEE
	DATE
	ACCOUNT_SEPARATOR
	CLEARED_INDICATOR
	IS_NEGATIVE
	COMMENT
	PRICE
	PRICE_CALC
	CURRENCY
)

type primitive struct {
	bytes.Buffer
	number int
	text   int
	other  map[rune]int
	last   rune
}

func (p *primitive) total() int {
	total := p.number + p.text
	for _, v := range p.other {
		total = total + v
	}
	return total
}

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

func isComment(ch rune) bool {
	return ch == ';'
}

func isCalculation(ch rune) bool {
	return ch == '(' || ch == ')' || ch == '+' || ch == '-'
}

func isText(ch rune) bool {
	return !isComment(ch) && !isAccountSeparator(ch) && (unicode.IsLetter(ch) || unicode.IsPunct(ch))
}

func isNegative(ch rune) bool {
	return ch == '-'
}
