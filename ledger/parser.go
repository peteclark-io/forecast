package ledger

import "errors"

type parser struct {
	lexer *lexer
}

func (p *parser) Parse() error {

	for {
		token, text := p.lexer.Scan()
		if token != NUMBER {
			p.fastForwardTo(CR)
			continue
		}

		p.parseDate()
	}
}

func (p *parser) fastForwardTo(token Token) {
	for {
		t, _ := p.lexer.Scan()
		if token == t {
			break
		}
	}
}

func (p *parser) parseDate() error {
	t, r := p.lexer.Scan()
	if t != DATE_SEPARATOR {
		return errors.New("Expected: / but found: " + r)
	}

}
