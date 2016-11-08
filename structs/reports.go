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
	Account           []string `json:"account"`
	Amount            float64  `json:"amount"`
	ExchangedAmount   float64  `json:"exchangedAmount"`
	Reported          bool     `json:"reported"`
	Currency          string   `json:"currency"`
	ExchangedCurrency string   `json:"exchangedCurrency"`
	Calculated        bool     `json:"calculated"`
	Virtual           bool     `json:"virtual"`
	Exchanged         bool     `json:"exchanged"`
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
	e.ExchangedAmount = 0
	e.ExchangedCurrency = ""
	e.Exchanged = false
}
