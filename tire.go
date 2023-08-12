package go_otzi

import (
	"fmt"
	"strings"
)

type ColumnNode struct {
	columns  map[string]struct{}
	requests []int // request IDs
}

type TableNode struct {
	children map[string]*ColumnNode
}

type DatabaseNode struct {
	children map[string]*TableNode
}

type ActionNode struct {
	children map[string]*DatabaseNode
}

type Trie struct {
	root *ActionNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &ActionNode{
			children: make(map[string]*DatabaseNode),
		},
	}
}

func (t *Trie) Insert(req SQLRequest) {
	currentAction := t.root

	// Insert action
	if _, ok := currentAction.children[req.Action]; !ok {
		currentAction.children[req.Action] = &DatabaseNode{
			children: make(map[string]*TableNode),
		}
	}
	currentDB := currentAction.children[req.Action]

	// Insert database
	if _, ok := currentDB.children[req.DB]; !ok {
		currentDB.children[req.DB] = &TableNode{
			children: make(map[string]*ColumnNode),
		}
	}
	currentTable := currentDB.children[req.DB]

	// Insert table
	if _, ok := currentTable.children[req.Table]; !ok {
		currentTable.children[req.Table] = &ColumnNode{
			columns: make(map[string]struct{}),
		}
	}

	currentColumn := currentTable.children[req.Table]

	// Insert columns
	for _, col := range req.Cols {
		currentColumn.columns[col] = struct{}{}
	}
	currentColumn.requests = append(currentColumn.requests, req.ID)
}

func (t *Trie) Merge() ([]string, map[string][]int) {
	results := []string{}
	reqMap := make(map[string][]int) // Map from merged SQL to its associated request IDs

	for action, dbNode := range t.root.children {
		for db, tableNode := range dbNode.children {
			for table, colNode := range tableNode.children {
				cols := make([]string, 0, len(colNode.columns))
				for col := range colNode.columns {
					cols = append(cols, col)
				}
				query := fmt.Sprintf("%s %s FROM %s.%s", action, strings.Join(cols, ", "), db, table)
				results = append(results, query)
				reqMap[query] = colNode.requests
			}
		}
	}
	return results, reqMap
}
