package ledger

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
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
		log.Println(posting)
	}

	return postings, nil
}

func (p *parser) parsePosting(token Token, value string) (structs.Posting, error) {
	posting := &structs.Posting{}
	if token != DATE {
		return *posting, errors.New("Why isn't this a date?!")
	}

	d, _ := time.Parse("2006/01/02", value)
	posting.Date = d

	var entries []structs.Entry
	entry := &structs.Entry{Amount: 1}
	last := token

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
			entry.Account = append(entry.Account, val)
		}

		if token == IS_NEGATIVE {
			entry.Amount = -1
		}

		if token == CURRENCY {
			entry.Currency = val
		}

		if token == PRICE {
			price, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return *posting, err
			}

			entry.Amount = entry.Amount * price
			entry.Reported = true
		}

		if token == CR && (last == PRICE || last == ACCOUNT) {
			entries = append(entries, *entry)
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
