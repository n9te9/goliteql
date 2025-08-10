package schema

func withBuiltin(s *Schema) *Schema {
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("Int"),
	})
	s.Indexes.ScalarIndex["Int"] = &ScalarDefinition{
		Name: []byte("Int"),
	}
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("Float"),
	})
	s.Indexes.ScalarIndex["Float"] = &ScalarDefinition{
		Name: []byte("Float"),
	}
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("String"),
	})
	s.Indexes.ScalarIndex["String"] = &ScalarDefinition{
		Name: []byte("String"),
	}
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("Boolean"),
	})
	s.Indexes.ScalarIndex["Boolean"] = &ScalarDefinition{
		Name: []byte("Boolean"),
	}
	s.Scalars = append(s.Scalars, &ScalarDefinition{
		Name: []byte("ID"),
	})
	s.Indexes.ScalarIndex["ID"] = &ScalarDefinition{
		Name: []byte("ID"),
	}

	return s
}
