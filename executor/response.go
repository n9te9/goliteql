package executor

import (
	"encoding/json"

	"github.com/lkeix/gg-executor/query"
)

func ExcludeSelectFields(resp map[string]json.RawMessage, selectSets []query.Selection) map[string]json.RawMessage {
	included := make(map[string]struct{}, len(selectSets))
	for _, sel := range selectSets {
		switch s := sel.(type) {
		case *query.Field:
			if _, ok := resp[string(s.Name)]; ok {
				included[string(s.Name)] = struct{}{}
			}
		}
	}

	for k := range resp {
		if _, ok := included[k]; !ok {
			delete(resp, k)
		}
	}

	return resp
}
