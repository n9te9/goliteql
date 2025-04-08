package executor

import "github.com/lkeix/gg-executor/query"

type Node struct {
	Name       []byte
	SelectSets []query.Selection
	Children   []*Node
}

func PlanExecution(selections []query.Selection) *Node {
	for _, sel := range selections {
		switch s := sel.(type) {
		case *query.Field:
			node := &Node{
				Name:       s.Name,
				SelectSets: s.Selections,
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
