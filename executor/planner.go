package executor

import (
	"github.com/n9te9/goliteql/query"
)

type Node struct {
	Name       []byte
	SelectSets []query.Selection
	Directives []*query.Directive
	Children   []*Node
}

func PlanExecution(selections []query.Selection) *Node {
	for _, sel := range selections {
		switch s := sel.(type) {
		case *query.Field:
			node := &Node{
				Name:       s.Name,
				SelectSets: s.Selections,
				Directives: s.Directives,
				Children:   make([]*Node, 0),
			}

			for _, child := range s.Selections {
				switch c := child.(type) {
				case *query.Field:
					node.Children = append(node.Children, digExecution(c))
				}
			}

			return node
		}
	}

	return nil
}

func digExecution(selectSet query.Selection) *Node {
	switch s := selectSet.(type) {
	case *query.Field:
		node := &Node{
			Name:       s.Name,
			SelectSets: s.Selections,
			Directives: s.Directives,
		}
		for _, child := range s.Selections {
			switch c := child.(type) {
			case *query.Field:
				node.Children = append(node.Children, digExecution(c))
			}
		}
		return node
	}

	return nil
}
