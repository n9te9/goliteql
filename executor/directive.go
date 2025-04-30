package executor

import (
	"encoding/json"

	"github.com/n9te9/goliteql/query"
)

func isIncludeDirective(dir *query.Directive) bool {
	if dir == nil {
		return false
	}

	if dir.Name != nil && string(dir.Name) == "include" {
		return true
	}

	return false
}

func IsIncluded(directives []*query.Directive, v json.RawMessage) bool {
	for _, dir := range directives {
		if isIncludeDirective(dir) {
			if len(dir.Arguments) != 1 {
				return false
			}

			if dir.Arguments[0].IsVariable {
				if string(dir.Arguments[0].Name) != "if" {
					return false
				}

				variables := make(map[string]any)
				if err := json.Unmarshal(v, &variables); err != nil {
					return false
				}

				flag, ok := variables[string(dir.Arguments[0].Name)].(bool)
				if !ok {
					return true
				}

				return flag
			}

			if string(dir.Arguments[0].Value) == "true" {
				return true
			}

			return false
		}
	}

	return true
}

func isSkipDirective(dir *query.Directive) bool {
	if dir == nil {
		return false
	}

	if dir.Name != nil && string(dir.Name) == "skip" {
		return true
	}

	return false
}

func IsSkipped(directives []*query.Directive, v json.RawMessage) bool {
	for _, dir := range directives {
		if isSkipDirective(dir) {
			if len(dir.Arguments) != 1 {
				return false
			}

			if dir.Arguments[0].IsVariable {
				if string(dir.Arguments[0].Name) != "if" {
					return true
				}

				variables := make(map[string]any)
				if err := json.Unmarshal(v, &variables); err != nil {
					return false
				}

				flag, ok := variables[string(dir.Arguments[0].Name)].(bool)
				if !ok {
					return true
				}

				return flag
			}

			return string(dir.Arguments[0].Value) == "true"
		}
	}

	return false
}
