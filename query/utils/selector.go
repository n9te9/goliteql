package utils

import "github.com/lkeix/gg-executor/query"

func ExtractSelectorName(op *query.Operation, operationName string) string {
	res := make([]string, 0, len(op.Selections))

	for _, sel := range op.Selections {
		switch s := sel.(type) {
		case *query.Field:
			res = append(res, string(s.Name))
		}
	}

	if len(res) == 1 {
		return res[0]
	}

	for _, s := range res {
		if s == operationName {
			return s
		}
	}

	return ""
}
