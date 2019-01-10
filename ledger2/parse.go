package ledger2

import (
	"bufio"
	"strings"
)

func computePrice(n *Node) {
	if n.Type != Price {
		return
	}

	r := bufio.NewReader(strings.NewReader(n.Value))
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return
		}

	}
}
