package executor

import (
	"github.com/n9te9/goliteql/query"
)

type Node struct {
	Name       []byte
	Directives []*query.Directive
	Arguments  []*query.Argument
	Type       string
	IsFragment bool
	Parent     *Node
	Children   []*Node
}

func PlanExecution(selections []query.Selection) []*Node {
	ret := make([]*Node, 0, len(selections))
	for _, sel := range selections {
		switch s := sel.(type) {
		case *query.Field:
			node := &Node{
				Name:       s.Name,
				Directives: s.Directives,
				Children:   make([]*Node, 0, len(s.Selections)),
				Arguments:  s.Arguments,
			}

			for _, child := range s.Selections {
				node.Children = append(node.Children, digExecution(nil, child))
			}

			ret = append(ret, node)
		}
	}

	return ret
}

func digExecution(parent *Node, selectSet query.Selection) *Node {
	switch s := selectSet.(type) {
	case *query.Field:
		node := &Node{
			Name:       s.Name,
			Directives: s.Directives,
			Arguments:  s.Arguments,
			Children:   make([]*Node, 0, len(s.Selections)),
			Parent:     parent,
		}
		for _, child := range s.Selections {
			switch c := child.(type) {
			case *query.Field:
				node.Children = append(node.Children, digExecution(node, c))
			case *query.InlineFragment:
				node.Parent.IsFragment = true
				node.Type = string(c.TypeCondition)
				node.Children = append(node.Children, digExecution(node, c))
			case *query.FragmentSpread:
				// Handle fragment spread
			}
		}
		return node
	case *query.InlineFragment:
		node := &Node{
			Directives: s.Directives,
			Type:       string(s.TypeCondition),
			Children:   make([]*Node, 0, len(s.Selections)),
			Parent:     parent,
		}
		for _, child := range s.Selections {
			switch c := child.(type) {
			case *query.Field:
				node.Children = append(node.Children, digExecution(node, c))
			case *query.InlineFragment:
				node.Parent.IsFragment = true
				node.Children = append(node.Children, digExecution(node, c))
			case *query.FragmentSpread:
				// Handle fragment spread
			}
		}
		return node
	}

	return nil
}
