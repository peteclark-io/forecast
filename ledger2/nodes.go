package ledger2

import "time"

func (l *Lexer) lexNode(expected Expectation, tk Token) (*Node, error) {
	switch expected.Node {
	case Date:
		return l.lexDateNode(tk)
	case Pending:
		fallthrough
	case Cleared:
		return l.lexPendingOrCleared(tk)
	case Payee:
		return l.lexPayee(tk)
	case Posting:
		return l.lexPosting(tk)
	case Price:
		return l.lexPrice(tk)
	case LineBreak:
		return l.lexLineBreak(tk)
	}

	return nil, syntaxErr
}

func (l *Lexer) lexPendingOrCleared(tk Token) (*Node, error) {
	if tk.Value == "*" {
		return &Node{Type: Cleared, Value: tk.Value}, nil
	} else if tk.Value == "!" {
		return &Node{Type: Pending, Value: tk.Value}, nil
	}
	return nil, syntaxErr
}

func (l *Lexer) lexDateNode(tk Token) (*Node, error) {
	_, err := time.Parse("2006/01/02", tk.Value)
	if err != nil {
		return nil, err
	}
	return &Node{Type: Date, Value: tk.Value}, nil
}

func (l *Lexer) lexPayee(tk Token) (*Node, error) {
	payee := tk.Value
	for {
		tk := l.t.Scan()
		if tk.Type == EOL || tk.Type == EOF || tk.Type == COMMENT {
			break
		}
		payee += tk.Value
	}

	return &Node{Type: Payee, Value: payee}, nil
}

func (l *Lexer) lexPosting(tk Token) (*Node, error) {
	posting := tk.Value
	for {
		tk := l.t.Scan()

		if tk.Type == WHITESPACE && len(tk.Value) == 1 {
			posting += tk.Value
			continue
		}

		if tk.Type == WHITESPACE || tk.Type == EOL || tk.Type == EOF || tk.Type == COMMENT {
			l.t.Undo()
			break
		}

		posting += tk.Value
	}

	return &Node{Type: Posting, Value: posting}, nil
}

func (l *Lexer) lexPrice(tk Token) (*Node, error) {
	price := tk.Value
	for {
		tk := l.t.Scan()
		if tk.Type == EOF || tk.Type == EOL || tk.Type == COMMENT {
			l.t.Undo()
			break
		}
		price += tk.Value
	}
	return &Node{Type: Price, Value: price}, nil
}

func (l *Lexer) lexLineBreak(tk Token) (*Node, error) {
	return &Node{Type: LineBreak, Value: tk.Value}, nil
}
