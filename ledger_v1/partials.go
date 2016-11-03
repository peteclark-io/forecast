package ledger_v1

import (
	"log"
	"time"
)

type partial struct {
	raw string
}

func (p partial) toDate() time.Time {
	log.Println(p.raw)
	d, _ := time.Parse("2006/01/02", p.raw)
	return d
}
