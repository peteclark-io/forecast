package ledger

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	p := NewParser(strings.NewReader(testPosting))
	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}
