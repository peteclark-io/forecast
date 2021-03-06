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

2016/11/15 ! Vodafone
    Expenses:Phone                                                  £45.40 ; ~
    Assets:Current:HSBC

2016/07/09 * Palm 2
    Expenses:Food:Breakfast                                          £4.30
    (Expenses:Food:Coffee)                                           £4.30
    Assets:Cash

2016/07/02 * Street Feast
    Expenses:Food:Dinner Out                                          (£14 + £8 + (£4.75 * 3))
    Assets:Cash                                                       (-£14 - £8)
    Assets:Current:HSBC

2016/09/17 * McDonalds
    Expenses:Food:Drunk                                             €2.58 @@ £2.19
    Assets:Current:Mondo

2016/09/18 * OktoberFest
    Expenses:Holiday:Oktoberfest                                    €122 @ £0.8496
    Assets:Luke                                                     €78  @ £0.8496
    Equity:Maths:Adjustments                                        £0.02 ; ^^ this is rounded incorrectly, so adjusting.
    Assets:Current:Mondo

2016/07/23 * Tesco
    Expenses:Food:Dinner                                            £1.92
    Expenses:Alcohol                                                £6.50
    Expenses:Food:Breakfast                                         (£2.45 + £2.10)
    Assets:Current:HSBC
`

func TestDateLex(t *testing.T) {
	l := lexer{stream: bufio.NewReader(strings.NewReader(testDate)), last: SOF}
	for {
		token, _ := l.Scan()
		//t.Log(token, txt)
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

func TestPriceCalc(t *testing.T) {
	l := lexer{stream: bufio.NewReader(strings.NewReader(testCalcPosting1)), last: SOF}
	for {
		token, _ := l.Scan()
		//t.Log(token, txt)
		if token == EOF {
			break
		}
	}
}
