package ledger

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	p := NewParser(strings.NewReader(testPosting))
	postings, _ := p.Parse()

	d, _ := json.MarshalIndent(postings, "", "  ")
	t.Log(string(d))
}
