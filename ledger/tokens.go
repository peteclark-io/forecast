package ledger

import (
	"bytes"
	"unicode"
)

type Token int

const (
	eof                       = rune(0)
	ILLEGAL             Token = iota // 1
	SOF                              // 2
	EOF                              // 3
	SPACE                            // 4
	BLOCK_SPACE                      // 5
	CR_BLOCK_SPACE                   // 6
	CR                               // 7
	ACCOUNT                          // 8
	PAYEE                            // 9
	DATE                             // 10
	ACCOUNT_SEPARATOR                // 11
	CLEARED_INDICATOR                // 12
	IS_NEGATIVE                      // 13
	COMMENT                          // 14
	PRICE                            // 15
	PRICE_CALC_BOUNDARY              // 16
	PRICE_OPERATOR                   // 17
	CURRENCY                         // 18
	VIRTUAL_ACCOUNT                  // 19
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

func isCalculationBoundary(ch rune) bool {
	return ch == '(' || ch == ')'
}

func isCalculationOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/'
}

func isText(ch rune) bool {
	return !isComment(ch) && !isAccountSeparator(ch) && (unicode.IsLetter(ch) || unicode.IsPunct(ch))
}

func isNegative(ch rune) bool {
	return ch == '-'
}
