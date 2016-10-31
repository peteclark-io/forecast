package structs

import "time"

type Report struct {
}

type Posting struct {
	Date    time.Time
	Cleared bool
	Payee   string
	Entries []Entry
}

type Entry struct {
	Account    string
	Amount     float64
	Unreported bool
	Currency   string
}
