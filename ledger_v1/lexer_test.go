package ledger_v1

import (
	"bufio"
	"strings"
	"testing"
)

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
    Assets:Savings:Two                                                £1401.54
    Assets:Savings:Thre                                               £21457.15
    Assets:Savings:Four                                               £18000
    Liabilities:Credit Card                                           -£646.21
    Equity:Opening Balances

2016/11/01 ! Vision Express
    Expenses:Contact Lenses                                         £17
    Assets:Current:Account
`

func TestScan(t *testing.T) {
	l := lexer{stream: bufio.NewReader(strings.NewReader(testPosting))}

	for {
		token, text := l.Scan()
		t.Log(text)
		if token == EOF {
			break
		}
	}
}
