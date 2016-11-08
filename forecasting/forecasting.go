package forecasting

import (
	"regexp"
	"strings"

	"github.com/peteclark-io/forecast/maths"
	"github.com/peteclark-io/forecast/structs"
)

func AverageByDay(filter *regexp.Regexp, postings []structs.Posting) map[string]results {
	data := make(map[string]results)

	for _, posting := range postings {
		for _, entry := range posting.Entries {
			account := join(entry.Account)
			if !filter.MatchString(account) {
				continue
			}

			r, ok := data[account]
			if !ok {
				r = results{
					Weekday: make([]*average, 7),
					Month:   make([]*average, 31),
				}
			}

			currentWeek := r.Weekday[posting.Date.Weekday()]
			currentMonth := r.Month[posting.Date.Day()-1]

			if currentWeek == nil {
				currentWeek = &average{}
				r.Weekday[posting.Date.Weekday()] = currentWeek
			}

			if currentMonth == nil {
				currentMonth = &average{}
				r.Month[posting.Date.Day()-1] = currentMonth
			}

			currentWeek.add(price(posting, entry))
			currentMonth.add(price(posting, entry))
			data[account] = r
		}
	}

	return data
}

type results struct {
	Weekday []*average `json:"weekday"`
	Month   []*average `json:"months"`
}

type average struct {
	Total   float64 `json:"total"`
	Count   int     `json:"count"`
	Average float64 `json:"avg"`
}

func (r *average) add(v float64) {
	r.Total = r.Total + v
	r.Count++
	r.Average = r.Total / float64(r.Count)
}

func price(posting structs.Posting, entry structs.Entry) float64 {
	if entry.Reported {
		return entry.Amount
	}

	return calcUnreported(posting)
}

func calcUnreported(posting structs.Posting) float64 {
	total := float64(0)
	for _, entry := range posting.Entries {
		if entry.Virtual {
			continue
		}

		if entry.Exchanged {
			total = total + entry.ExchangedAmount
			continue
		}

		total = total + entry.Amount
	}

	return maths.Round(total, 2) * -1
}

func join(account []string) string {
	return strings.Join(account, ":")
}
