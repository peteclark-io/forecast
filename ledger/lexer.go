package ledger

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"unicode"
)

type lexer struct {
	stream *bufio.Reader
	last   Token
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

func (l *lexer) unreadCharacters(chars int) {
	for i := 0; i < chars; i++ {
		l.stream.UnreadRune()
	}
}

func (l *lexer) Scan() (Token, string) {
	ch := l.read()

	if ch == eof {
		return EOF, string(ch)
	}

	if (l.last == COMMENT) && !isNewLine(ch) {
		return COMMENT, string(ch)
	}

	if isNumber(ch) && (l.last == CR || l.last == SOF) {
		l.unread()
		date, err := l.readDate(ch)
		if err != nil {
			log.Println(err)
			return ILLEGAL, string(ch)
		}

		l.last = DATE
		return DATE, date
	}

	if isWhitespace(ch) {
		l.unread()
		ws := l.readAllOfType(isWhitespace, nil)

		if len(ws) > 1 && l.last == CR {
			l.last = CR_BLOCK_SPACE
			return CR_BLOCK_SPACE, ws
		} else if len(ws) > 1 {
			l.last = BLOCK_SPACE
			return BLOCK_SPACE, ws
		}

		return SPACE, ws
	}

	if isCleared(ch) && l.last == DATE {
		l.last = CLEARED_INDICATOR
		return CLEARED_INDICATOR, string(ch)
	}

	if (isText(ch) || isNumber(ch)) && l.last == CLEARED_INDICATOR {
		l.unread()
		txt := l.readAllOfType(func(ch rune) bool {
			return isText(ch) || isWhitespace(ch) || isNumber(ch)
		}, nil)
		return PAYEE, txt
	}

	if isText(ch) && (l.last == CR_BLOCK_SPACE || l.last == ACCOUNT_SEPARATOR) {
		l.unread()
		txt := l.readAccount(ch)
		l.last = ACCOUNT
		return ACCOUNT, txt
	}

	if isAccountSeparator(ch) && l.last == ACCOUNT {
		l.last = ACCOUNT_SEPARATOR
		return ACCOUNT_SEPARATOR, string(ch)
	}

	if isCalculationBoundary(ch) && (l.last == BLOCK_SPACE || l.last == PRICE || l.last == PRICE_OPERATOR) {
		return PRICE_CALC_BOUNDARY, string(ch)
	}

	if isNegative(ch) && (l.last == BLOCK_SPACE) {
		l.last = IS_NEGATIVE
		return IS_NEGATIVE, string(ch)
	}

	if isNegative(ch) && (l.last == BLOCK_SPACE || l.last == PRICE_OPERATOR) {
		l.last = IS_NEGATIVE
		return IS_NEGATIVE, string(ch)
	}

	if unicode.IsSymbol(ch) && (l.last == BLOCK_SPACE || l.last == IS_NEGATIVE || l.last == PRICE_CALC_BOUNDARY || l.last == PRICE_OPERATOR) {
		l.unread()
		txt := l.readAllOfType(unicode.IsSymbol, nil)
		l.last = CURRENCY
		return CURRENCY, txt
	}

	if isNumber(ch) && (l.last == CURRENCY || l.last == PRICE_OPERATOR) {
		l.unread()
		price := l.readAllOfType(isNumber, nil)
		l.last = PRICE
		return PRICE, price
	}

	if isCalculationOperator(ch) && l.last == PRICE {
		l.last = PRICE_OPERATOR
		return PRICE_OPERATOR, string(ch)
	}

	if isNewLine(ch) {
		l.last = CR
		return CR, string(ch)
	}

	if isComment(ch) {
		l.last = COMMENT
		return COMMENT, string(ch)
	}

	l.last = ILLEGAL
	return ILLEGAL, string(ch)
}

func (l *lexer) readAccount(ch rune) string {
	p := &primitive{text: 1, last: ch}
	p.WriteRune(l.read())

	for {
		ch := l.read()

		if ch == eof {
			break
		}

		if isWhitespace(ch) && !isWhitespace(p.last) {
			p.last = ch
			p.WriteRune(ch)
			continue
		}

		if !isText(ch) {
			l.unread()
			break
		}

		p.last = ch
		p.WriteRune(ch)
	}

	return p.String()
}

func (l *lexer) readAllOfType(check func(ch rune) bool, ignore func(ch rune) bool) string {
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	for {
		ch := l.read()

		if ch == eof {
			break
		}

		if ignore != nil && ignore(ch) {
			continue
		}

		if !check(ch) {
			l.unread()
			break
		}

		buf.WriteRune(ch)
	}

	return buf.String()
}

func (l *lexer) readDate(ch rune) (string, error) {
	p := &primitive{number: 1, other: make(map[rune]int), last: ch}
	p.WriteRune(l.read())

	for {
		if p.total() == 10 {
			break
		}

		ch := l.read()

		if ch == eof {
			break
		}

		if isNumber(ch) {
			p.number = p.number + 1
			p.last = ch
		}

		if isDateSeparator(ch) {
			p.other[ch] = p.other[ch] + 1
			p.last = ch
		}

		if !isNumber(ch) && !isDateSeparator(ch) {
			l.unreadCharacters(p.total())
			return "", errors.New("Not a date")
		}

		p.WriteRune(ch)
	}

	return p.String(), nil
}
