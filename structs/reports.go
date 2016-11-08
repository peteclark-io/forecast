package structs

import "time"

type Report struct {
}

type Posting struct {
	Date    time.Time `json:"date"`
	Cleared bool      `json:"cleared"`
	Payee   string    `json:"payee"`
	Entries []Entry   `json:"entries"`
}

type Entry struct {
	Account    []string `json:"account"`
	Amount     float64  `json:"amount"`
	Reported   bool     `json:"reported"`
	Currency   string   `json:"currency"`
	Calculated bool     `json:"calculated"`
	Virtual    bool     `json:"virtual"`
}

func (e *Entry) IsComplete() bool {
	return len(e.Account) > 0
}

func (e *Entry) Reset() {
	e.Account = []string{}
	e.Amount = 0
	e.Reported = false
	e.Calculated = false
	e.Virtual = false
	e.Currency = ""
}
