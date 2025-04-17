package executor

import (
	"encoding/json"
	"slices"

	"github.com/lkeix/gg-executor/query"
)

func ExcludeSelectFields(resp map[string]json.RawMessage, selectSets []query.Selection) map[string]json.RawMessage {
	included := make([]string, 0)
	for _, sel := range selectSets {
		switch s := sel.(type) {
		case *query.Field:
			if _, ok := resp[string(s.Name)]; ok {
				included = append(included, string(s.Name))
			}
		}
	}

	for k := range resp {
		if !slices.Contains(included, k) {
			delete(resp, k)
		}
	}

	return resp
}
