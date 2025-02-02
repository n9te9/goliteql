package schema

import (
	"fmt"
	"strconv"
)


var typesValidator = map[string]func([]byte) error{
	"Int": func(value []byte) error {
		if _, err := strconv.Atoi(string(value)); err != nil {
			return fmt.Errorf("expected integer value, got %s", value)
		}

		return nil
	},
	"Float": func(value []byte) error {
		if _, err := strconv.ParseFloat(string(value), 64); err != nil {
			return fmt.Errorf("expected float value, got %s", value)
		}

		return nil
	},
	"String": func(value []byte) error {
		if value[0] != '"' || value[len(value)-1] != '"' {
			return fmt.Errorf("expected String but got %s", value)
		}
		return nil
	},
	"Boolean": func(value []byte) error {
		if _, err := strconv.ParseBool(string(value)); err != nil {
			return fmt.Errorf("expected boolean value, got %s", value)
		}

		return nil
	},
	"ID": func(value []byte) error {
		if _, err := strconv.Atoi(string(value)); err != nil || value[0] != '"' || value[len(value)-1] != '"' {
			return fmt.Errorf("expected ID but got %s", value)
		}
		return nil
	},
}

type ArgumentDefinition struct {
	Name []byte
	Default []byte
	Type *FieldType
}

func (a *ArgumentDefinition) ValidateValueType(value []byte) error {
	if err := typesValidator[string(a.Type.Name)](value); err != nil {
		return fmt.Errorf("error validating value for argument %s: %w", a.Name, err)
	}

	return nil
}

type ArgumentDefinitions []*ArgumentDefinition

func (a ArgumentDefinitions) RequiredArguments() map[*ArgumentDefinition]struct{} {
	res := make(map[*ArgumentDefinition]struct{})
	for _, arg := range a {
		if !arg.Type.Nullable {
			res[arg] = struct{}{}
		}
	}

	return res
}