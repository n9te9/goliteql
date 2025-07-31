package executor

import (
	"encoding/json"

	"github.com/n9te9/goliteql/query"
)

type Directives []*query.Directive

func (d Directives) ShouldInclude(variables map[string]json.RawMessage) bool {
	return isIncluded(d, variables) && !isSkipped(d, variables)
}

func isIncludeDirective(dir *query.Directive) bool {
	if dir == nil {
		return false
	}

	if dir.Name != nil && string(dir.Name) == "include" {
		return true
	}

	return false
}

func isIncluded(directives []*query.Directive, variables map[string]json.RawMessage) bool {
	for _, dir := range directives {
		if isIncludeDirective(dir) {
			if len(dir.Arguments) != 1 {
				return false
			}

			if dir.Arguments[0].IsVariable {
				if string(dir.Arguments[0].Name) != "if" {
					return false
				}

				flag, ok := variables[string(dir.Arguments[0].Name)]
				if !ok {
					return true
				}

				return string(flag) == "true"
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

func isSkipped(directives []*query.Directive, variables map[string]json.RawMessage) bool {
	for _, dir := range directives {
		if isSkipDirective(dir) {
			if len(dir.Arguments) != 1 {
				return false
			}

			if dir.Arguments[0].IsVariable {
				if string(dir.Arguments[0].Name) != "if" {
					return false
				}

				flag, ok := variables[string(dir.Arguments[0].Name)]
				if !ok {
					return true
				}

				return string(flag) == "true"
			}

			if string(dir.Arguments[0].Value) == "true" {
				return true
			}

			return false
		}
	}

	return false
}
