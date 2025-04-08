package utils

import (
	"encoding/json"
	"fmt"

	"github.com/lkeix/gg-executor/query"
)

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

func ExtractSelectorArgs(op *query.Operation, operationName string) []*query.Argument {
	if op == nil {
		return nil
	}

	for _, sel := range op.Selections {
		switch s := sel.(type) {
		case *query.Field:
			if string(s.Name) == operationName {
				return s.Arguments
			}
		}
	}

	return nil
}

func ExtractExecuteSelector(op *query.Operation, operationName string) []query.Selection {
	if op == nil {
		return nil
	}

	for _, sel := range op.Selections {
		switch s := sel.(type) {
		case *query.Field:
			if string(s.Name) == operationName {
				return []query.Selection{s}
			}
		}
	}

	if len(op.Selections) == 1 {
		switch s := op.Selections[0].(type) {
		case *query.Field:
			return []query.Selection{s}
		}
	}

	return nil
}

func ConvRequestBodyFromVariables(variables json.RawMessage, args []*query.Argument) ([]byte, error) {
	mp := make(map[string]any)
	ret := make(map[string]any)

	if err := json.Unmarshal(variables, &mp); err != nil {
		return nil, err
	}

	for i, arg := range args {
		for k, v := range mp {
			if string(arg.Name) == k {
				ret[fmt.Sprintf("arg%d", i)] = v
			}
		}
	}

	return json.Marshal(ret)
}
