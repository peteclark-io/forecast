package ledger2

import (
	"bufio"
	"unicode"
)

type Tokenizer struct {
	stream    *bufio.Reader
	undo      bool
	lastToken Token
}

type Token struct {
	Type  TokenType
	Value string
}

type TokenType int

const (
	eof                  = rune(0)
	EOF        TokenType = iota // 1
	WHITESPACE                  // 2
	STATEMENT                   // 3
	EOL                         // 4
	COMMENT                     // 5
)

func (t *Tokenizer) Scan() (tk Token) {
	if t.undo {
		t.undo = false
		return t.lastToken
	}

	ch := t.readRune()

	defer func() {
		t.lastToken = tk
	}()

	if ch == eof {
		return Token{Type: EOF, Value: ""}
	}

	if isWhitespace(ch) {
		chs := t.readUntil(isNotWhitespace)
		chs = append([]rune{ch}, chs...)
		return Token{Type: WHITESPACE, Value: string(chs)}
	}

	if isNewLine(ch) {
		return Token{Type: EOL, Value: string(ch)}
	}

	if isComment(ch) {
		chs := t.readUntil(isNewLine)
		chs = append([]rune{ch}, chs...)
		return Token{Type: COMMENT, Value: string(chs)}
	}

	chs := t.readUntil(isWhitespaceOrComment)
	chs = append([]rune{ch}, chs...)
	return Token{Type: STATEMENT, Value: string(chs)}
}

func (t *Tokenizer) readRune() rune {
	ch, _, err := t.stream.ReadRune()

	if err != nil {
		return eof
	}
	return ch
}

func (t *Tokenizer) Undo() {
	if t.undo {
		panic("please undo just once")
	}
	t.undo = true
}

func (t *Tokenizer) readUntil(predicate func(r rune) bool) []rune {
	b := make([]rune, 0)
	for {
		ch := t.readRune()
		if predicate(ch) || ch == eof {
			t.stream.UnreadRune()
			break
		}

		b = append(b, ch)
	}
	return b
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\v' || ch == '\f'
}

func isComment(ch rune) bool {
	return ch == ';'
}

func isNewLine(ch rune) bool {
	return !isWhitespace(ch) && unicode.IsSpace(ch)
}

func isNotWhitespace(ch rune) bool {
	return !isWhitespace(ch) && !isNewLine(ch)
}

func isWhitespaceOrComment(ch rune) bool {
	return isWhitespace(ch) || isNewLine(ch) || isComment(ch)
}
