package forecasting

import (
	"strings"

	"github.com/peteclark-io/forecast/maths"
	"github.com/peteclark-io/forecast/structs"
)

func Balance(account string, postings []structs.Posting) (float64, error) {
	balance := float64(0)

	for _, posting := range postings {
		for _, entry := range posting.Entries {
			if entry.Virtual {
				continue
			}

			if strings.Join(entry.Account, ":") == account {
				amount := entry.Amount
				if !entry.Reported {
					amount = calcUnreported(posting)
				}

				balance = maths.Round(balance+amount, 2)
			}
		}
	}
	return balance, nil
}
