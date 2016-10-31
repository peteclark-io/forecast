package ledger

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
    Assets:Current:HSBC                                               £4.73
    Assets:Cash                                                       £62.50
    Assets:Savings:Flex                                               £11648.81
    Assets:Savings:Save To Buy ISA                                    £1403.86
    Assets:Savings:ISA                                                £21457.15
    Assets:Savings:Premium Bonds                                      £18000
    Liabilities:HSBC Credit Card                                      -£646.21
    Equity:Opening Balances
`

func TestScan(t *testing.T) {
	l := lexer{stream: bufio.NewReader(strings.NewReader(testPosting))}

	for {
		token, result := l.Scan()
		t.Log(token, result)
		if token == EOF {
			break
		}
	}
}
