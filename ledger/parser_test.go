package ledger

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testCalcPosting1 = `
2016/07/23 * Margate
    Assets:Cash                                                     £30
    Expenses:Fun                                                    ((£50 - £12.70) / 2)
    Expenses:Food:Dinner Out                                        (£12.70 / 2)
    Expenses:Food:Snacks                                            (£5.70 / 2)
    Assets:Ducky                                                    (((£50 - £12.70) / 2) + (£12.70 / 2) + (£5.70 / 2))
    Assets:Current:HSBC
`

const testCalcPosting2 = `
2016/11/01 ! Vision Express
    Expenses:Contact Lenses                                         (£17 + -£10 - £1)
    Assets:Current:Account ; what up doggg
 `

func TestParse(t *testing.T) {
	p := NewParser(strings.NewReader(testCalcPosting1))
	postings, _ := p.Parse()

	//d, _ := json.MarshalIndent(postings, "", "  ")
	//t.Log(string(d))

	assert.Len(t, postings, 1)
	assert.Equal(t, postings[0].Cleared, true)
	assert.Equal(t, postings[0].Payee, "Margate")

	for _, entry := range postings[0].Entries {
		switch strings.Join(entry.Account, ":") {
		case "Assets:Cash":
			assert.Equal(t, entry.Amount, float64(30))
			assert.False(t, entry.Calculated)
			assert.True(t, entry.Reported)
			break
		case "Expenses:Fun":
			assert.Equal(t, entry.Amount, 18.65)
			assert.True(t, entry.Calculated)
			assert.True(t, entry.Reported)
			break
		case "Expenses:Food:Dinner Out":
			assert.Equal(t, entry.Amount, 6.35)
			assert.True(t, entry.Calculated)
			assert.True(t, entry.Reported)
			break
		case "Expenses:Food:Snacks":
			assert.Equal(t, entry.Amount, 2.85)
			assert.True(t, entry.Calculated)
			assert.True(t, entry.Reported)
			break
		case "Assets:Ducky":
			assert.Equal(t, entry.Amount, 27.85)
			assert.True(t, entry.Calculated)
			assert.True(t, entry.Reported)
			break
		case "Assets:Current:HSBC":
			assert.Equal(t, entry.Amount, float64(0))
			assert.False(t, entry.Reported)
			assert.False(t, entry.Calculated)
			break
		}
	}
}

func TestParsePosting(t *testing.T) {
	p := NewParser(strings.NewReader(testPosting))
	p.Parse()

	//d, _ := json.MarshalIndent(postings, "", "  ")
	//t.Log(string(d))
}
