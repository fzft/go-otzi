package go_otzi

import (
	"fmt"
	"testing"
)

func TestTriePrint(t *testing.T) {

	trie := NewTrie()

	req1 := SQLRequest{
		ID:     1,
		Action: "SELECT",
		DB:     "mydb",
		Table:  "users",
		Cols:   []string{"name", "age"},
	}
	req2 := SQLRequest{
		ID:     2,
		Action: "SELECT",
		DB:     "mydb",
		Table:  "users",
		Cols:   []string{"email"},
	}
	req3 := SQLRequest{
		ID:     3,
		Action: "SELECT",
		DB:     "mydb",
		Table:  "orders",
		Cols:   []string{"order_id", "product"},
	}

	trie.Insert(req1)
	trie.Insert(req2)
	trie.Insert(req3)

	mergedQueries, reqMap := trie.Merge()

	for _, query := range mergedQueries {
		fmt.Println(query)
	}

	for req, query := range reqMap {
		fmt.Println(req, query)
	}

}
