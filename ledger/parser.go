package ledger

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/peteclark-io/forecast/structs"
)

type parser struct {
	lexer
}

type Parser interface {
	Parse() ([]structs.Posting, error)
}

func NewParser(input io.Reader) Parser {
	return &parser{lexer{bufio.NewReader(input), SOF}}
}

func (p *parser) Parse() ([]structs.Posting, error) {
	var postings []structs.Posting

	for {
		token, val := p.Scan()
		if token == EOF {
			break
		}

		if token != DATE {
			continue
		}

		posting, err := p.parsePosting(token, val)
		if err != nil {
			continue
		}

		postings = append(postings, posting)
	}

	return postings, nil
}

func (p *parser) parseCalculation(token Token, value string) (string, float64, error) {
	var price float64

	var currency *string
	var operator *string

	for {
		token, val := p.Scan()
		if token == EOF {
			break
		}

		if token == PRICE && operator == nil {
			p, err := strconv.ParseFloat(val, 64)
			price = p
			if err != nil {
				return "", 0, err
			}
		}

		if token == PRICE && operator != nil {
			p, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return "", 0, err
			}
			price = doCalc(price, p, *operator)
			operator = nil
		}

		if token == CURRENCY && currency == nil {
			currency = &val
		}

		if token == PRICE_OPERATOR {
			operator = &val
		}

		if token == PRICE_CALC_BOUNDARY && val == "(" {
			c, p, err := p.parseCalculation(token, value)
			if err != nil {
				return "", 0, err
			}

			currency = &c

			if operator != nil {
				price = doCalc(price, p, *operator)
			} else {
				price = doCalc(price, p, "+")
			}

			operator = nil
		}

		if token == PRICE_CALC_BOUNDARY && val == ")" {
			break
		}
	}

	return *currency, price, nil
}

func doCalc(p1, p2 float64, operator string) float64 {
	switch operator {
	case "+":
		p1 = p1 + p2
		break
	case "-":
		p1 = p1 - p2
		break
	case "/":
		p1 = p1 / p2
		break
	case "*":
		p1 = p1 * p2
		break
	}

	return Round(p1, 0.5, 2)
}

func (p *parser) parsePosting(token Token, value string) (structs.Posting, error) {
	posting := &structs.Posting{}
	if token != DATE {
		return *posting, errors.New("Why isn't this a date?!")
	}

	d, _ := time.Parse("2006/01/02", value)
	posting.Date = d

	var entries []structs.Entry
	entry := &structs.Entry{}
	last := token

	negative := float64(1)

	crs := 0
	for {
		token, val := p.Scan()
		if token == EOF {
			break
		}

		if token != CR {
			crs = 0
		}

		if token == CLEARED_INDICATOR {
			posting.Cleared = (val == "*")
		}

		if token == PAYEE {
			posting.Payee = val
		}

		if token == ACCOUNT {
			entry.Account = append(entry.Account, strings.TrimSpace(val))
		}

		if token == VIRTUAL_ACCOUNT {
			entry.Account = append(entry.Account, strings.TrimSpace(val))
			entry.Virtual = true
		}

		if token == PRICE_CALC_BOUNDARY {
			currency, price, err := p.parseCalculation(token, value)
			if err != nil {
				return *posting, nil
			}

			entry.Currency = currency
			entry.Amount = price
			entry.Calculated = true
			entry.Reported = true
			entries = append(entries, *entry)
			negative = 1
			entry.Reset()
		}

		if token == IS_NEGATIVE {
			negative = -1
		}

		if token == CURRENCY {
			entry.Currency = val
		}

		if token == PRICE {
			price, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return *posting, err
			}

			entry.Amount = negative * price
			entry.Reported = true
		}

		if token == CR && (last == PRICE || last == ACCOUNT || last == COMMENT) {
			entries = append(entries, *entry)
			negative = 1
			entry.Reset()
		}

		if token == CR {
			crs++
			if crs == 2 {
				break
			}
		}

		last = token
	}

	posting.Entries = entries

	return *posting, nil
}
