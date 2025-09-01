package schema

func WithBuiltin(s *Schema) *Schema {
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("Int"),
	})
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("Float"),
	})
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("String"),
	})
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("Boolean"),
	})
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("ID"),
	})

	return s
}
