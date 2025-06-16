package executor

import (
	"encoding/json"
	"errors"

	"github.com/n9te9/goliteql/query"
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

type GraphQLError struct {
	Message    string         `json:"message"`
	Path       []string       `json:"path,omitempty"`
	Extensions map[string]any `json:"extensions,omitempty"`
}

func (e GraphQLError) Error() string {
	return e.Message
}

type GraphQLResponse struct {
	Data   map[string]any `json:"data"`
	Errors []error        `json:"errors,omitempty"`
}

func MatchGraphQLResponse[T map[string]json.RawMessage | json.RawMessage | any](resp map[string]T) error {
	if _, ok := resp["errors"]; ok {
		var gqlErrors []GraphQLError
		if err := json.Unmarshal(any(resp["errors"]).(json.RawMessage), &gqlErrors); err != nil {
			return errors.New("failed to unmarshal GraphQL Error format\nuse executor.GraphQLError to unmarshal")
		}
	}

	if _, ok := resp["data"]; !ok {
		return errors.New("missing data field in GraphQL response\nuse executor.GraphQLResponse to unmarshal")
	}

	return nil
}
