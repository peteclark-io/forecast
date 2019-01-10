package ledger2

import (
	"errors"
)

var syntaxErr = errors.New("Unrecognised syntax")

var StartNode = &Node{Type: SOF, Value: ""}

type Statement struct {
	nodes []*Node
}

type NodeType int

type Node struct {
	Type  NodeType
	Value string
}

const (
	SOF       NodeType = iota // 0
	Date                      // 1
	Cleared                   // 2
	Pending                   // 3
	Posting                   // 4
	Payee                     // 5
	Price                     // 6
	LineBreak                 // 7
)

type Expectation struct {
	Node  NodeType
	Token TokenType
}

func (n NodeType) IgnoreProceedingTokens() []TokenType {
	switch n {
	case Date:
		return []TokenType{WHITESPACE}
	case Cleared:
		fallthrough
	case Pending:
		return []TokenType{WHITESPACE}
	case Payee:
		return []TokenType{EOL, WHITESPACE, COMMENT}
	case Posting:
		return []TokenType{WHITESPACE}
	case Price:
		return []TokenType{COMMENT}
	case LineBreak:
		return []TokenType{WHITESPACE}
	}
	return []TokenType{}
}

func (n NodeType) ExpectTokens() []Expectation {
	switch n {
	case SOF:
		return []Expectation{Expectation{Node: Date, Token: STATEMENT}}
	case Date:
		return []Expectation{Expectation{Node: Cleared, Token: STATEMENT}}
	case Cleared:
		fallthrough
	case Pending:
		return []Expectation{Expectation{Node: Payee, Token: STATEMENT}}
	case Payee:
		return []Expectation{Expectation{Node: Posting, Token: STATEMENT}}
	case Posting:
		return []Expectation{Expectation{Node: Price, Token: STATEMENT}, Expectation{Node: LineBreak, Token: EOL}, Expectation{Node: LineBreak, Token: COMMENT}}
	case Price:
		return []Expectation{Expectation{Node: LineBreak, Token: EOL}, Expectation{Node: LineBreak, Token: COMMENT}}
	case LineBreak:
		return []Expectation{Expectation{Node: Posting, Token: STATEMENT}}
	}
	return []Expectation{}
}

type Lexer struct {
	t         *Tokenizer
	lastToken Token
	lastNode  *Node
}

func (l *Lexer) Lex() ([]*Statement, error) {
	statements := make([]*Statement, 0)

eof:
	for {
		tk := l.t.Scan()

		switch tk.Type {
		case EOF:
			break eof
		case EOL:
			continue
		case STATEMENT:
			s, err := l.lexStatement(tk)
			if err != nil {
				continue
			}
			statements = append(statements, s)
			continue
		case WHITESPACE:
			continue
		}
	}

	return statements, nil
}

func (l *Lexer) lexStatement(tk Token) (*Statement, error) {
	statement := &Statement{nodes: make([]*Node, 0)}
	l.lastNode = StartNode

statement:
	for {
		if tk.Type == EOF {
			break statement
		}
		if tk.Type == EOL && l.lastToken.Type == EOL {
			break statement
		}

		if l.syntaxCheck(tk) {
			l.lastToken = tk
			tk = l.t.Scan()
			continue statement
		}

		for _, expectation := range l.lastNode.Type.ExpectTokens() {
			if expectation.Token == tk.Type {
				n, err := l.lexNode(expectation, tk)
				if err != nil {
					return nil, err
				}
				statement.nodes = append(statement.nodes, n)
				l.lastNode = n
			}
		}

		l.lastToken = tk
		tk = l.t.Scan()
	}

	return statement, nil
}

func (l *Lexer) syntaxCheck(tk Token) bool {
	types := l.lastNode.Type.IgnoreProceedingTokens()
	if contains(tk.Type, types) {
		return true
	}
	return false
}

func contains(t TokenType, types []TokenType) bool {
	for _, v := range types {
		if t == v {
			return true
		}
	}
	return false
}
