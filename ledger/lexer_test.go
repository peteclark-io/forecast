package ledger

import (
	"bufio"
	"strings"
	"testing"
)

const testDate = `2016/02/08 * I am a really long Payee
`

const testPosting = `
commodity £
    note GBP
    format £1,000.00
    nomarket
    default

2015/12/31 * Opening Balances
    Assets:Current:Account                                            £4.73
    Assets:Cash                                                       £62.50
    Assets:Savings:One                                                £171648.81
    Assets:Savings:Two                                                €1401.54
    Assets:Savings:Thre                                               £21457.15
    Assets:Savings:Four                                               -£18000
    Liabilities:Credit Card                                           -£646.21
    Equity:Opening Balances

2016/11/01 ! Vision Express
    Expenses:Contact Lenses                                         £17
    Assets:Current:Account ; what up doggg
`

func TestDateLex(t *testing.T) {
	l := lexer{stream: bufio.NewReader(strings.NewReader(testDate)), last: SOF}
	for {
		token, txt := l.Scan()
		t.Log(token, txt)
		if token == EOF {
			break
		}
	}
}

func TestFullPosting(t *testing.T) {
	l := lexer{stream: bufio.NewReader(strings.NewReader(testPosting)), last: SOF}
	for {
		token, txt := l.Scan()
		t.Log(token, txt)
		if token == EOF {
			break
		}
	}
}
