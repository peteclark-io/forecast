package balance

import (
	"log"
	"strings"

	"github.com/peteclark-io/forecast/ledger"
	"github.com/peteclark-io/forecast/structs"
)

func calcUnreported(posting structs.Posting) float64 {
	total := float64(0)
	for _, entry := range posting.Entries {
		total = total + entry.Amount
	}

	return total * -1
}

func Balance(account string, postings []structs.Posting) (float64, error) {
	balance := float64(0)

	for _, posting := range postings {
		for _, entry := range posting.Entries {
			if strings.Join(entry.Account, ":") == account {
				amount := entry.Amount
				if !entry.Reported {
					amount = calcUnreported(posting)
				}
				log.Println(posting)
				balance = ledger.Round(balance+amount, 0.5, 2)
				log.Println(amount, balance)
			}
		}
	}
	return balance, nil
}
