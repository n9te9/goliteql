package executor

import "encoding/json"

type Nullable interface {
	MarshalJSON() ([]byte, error)
}

type nullable struct {
	data any
}

func NewNullable(value any) *nullable {
	return &nullable{data: value}
}

func (a *nullable) MarshalJSON() ([]byte, error) {
	if a == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(a.data)
}
