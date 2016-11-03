package ledger_v1

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/peteclark-io/forecast/structs"
)

type parser struct {
	lexer *lexer
}

type LedgerParser interface {
	Parse() error
}

func NewParser(reader io.Reader) LedgerParser {
	return &parser{lexer: &lexer{bufio.NewReader(reader)}}
}

func (p *parser) Parse() error {
	for {
		token, text := p.lexer.Scan()
		if token == EOF {
			break
		}

		if token == WS || token == CR {
			continue
		}

		if token != NUMBER {
			p.fastForwardTo(func(token Token) bool {
				return token == CR
			})
			continue
		}

		date := &partial{text}
		err := p.parseDate(date)
		if err != nil {
			return err
		}

		_, err = p.cleared()
		if err != nil {
			return err
		}

		payee := &partial{}
		p.readPayee(payee)

		for {
			entry, err := p.parseEntry()
			if err != nil {
				break
			}

			log.Println(entry)
		}
	}

	return nil
}

func (p *parser) fastForwardTo(check func(token Token) bool) (Token, string) {
	for {
		t, r := p.lexer.Scan()
		if check(t) || t == EOF {
			return t, r
		}
	}
}

func (p *parser) parseEntry() (structs.Entry, error) {
	var entry structs.Entry
	account, err := p.parseAccount()
	if err != nil {
		return entry, err
	}

	t, r := p.lexer.Scan()

	if t == CR || t == EOF {
		return structs.Entry{
			Account:    account,
			Unreported: true,
		}, nil
	}

	currency := &partial{r}
	amount := r
	if t != NUMBER {
		amount = p.parseCurrency(currency)
	} else {
		currency.raw = ""
	}

	t, r = p.lexer.Scan()
	amount64, _ := strconv.ParseFloat(amount, 64)
	return structs.Entry{
		Account:    account,
		Currency:   currency.raw,
		Unreported: false,
		Amount:     amount64,
	}, nil
}

func (p *parser) parseCurrency(currency *partial) string {
	for {
		t, r := p.lexer.Scan()
		if t == NUMBER {
			return r
		}
		currency.raw = currency.raw + r
	}
}

func (p *parser) parseAccount() (string, error) {
	t, r := p.fastForwardTo(func(token Token) bool {
		return token != WS
	})

	if t == EOF || t == CR {
		return "", errors.New("End of posting")
	}

	account := partial{r}
	for {
		t, r = p.lexer.Scan()
		if t == EOF || t == CR {
			break
		}

		if t == WS && len(r) > 2 {
			break
		}

		account.raw = account.raw + r
	}
	return account.raw, nil
}

func (p *parser) cleared() (bool, error) {
	for {
		t, r := p.lexer.Scan()
		if t == WS {
			continue
		}

		if t != CLEARED_INDICATOR {
			return false, errors.New("Expected cleared indicator! Found: " + r)
		}

		return r == "*", nil
	}
}

func (p *parser) readPayee(payee *partial) {
	for {
		t, r := p.lexer.Scan()
		if t == CR {
			break
		}

		payee.raw = payee.raw + r
	}
	payee.raw = strings.TrimSpace(payee.raw)
}

func (p *parser) parseDate(date *partial) error {
	t, r := p.lexer.Scan()
	if t != DATE_SEPARATOR {
		return errors.New("Expected: / but found: " + r)
	}

	date.raw = date.raw + r

	t, r = p.lexer.Scan()
	if t != NUMBER {
		return errors.New("Expected: NUMBER but found: " + r)
	}

	date.raw = date.raw + r

	t, r = p.lexer.Scan()
	if t != DATE_SEPARATOR {
		return errors.New("Expected: / but found: " + r)
	}

	date.raw = date.raw + r

	t, r = p.lexer.Scan()
	if t != NUMBER {
		return errors.New("Expected: NUMBER but found: " + r)
	}

	date.raw = date.raw + r
	return nil
}
