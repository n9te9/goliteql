package executor

import (
	"github.com/n9te9/goliteql/query"
)

type Node struct {
	Name       string
	Directives []*query.Directive
	Arguments  []*query.Argument
	Type       string
	Parent     *Node
	Children   []*Node
}

func (n *Node) HasFragment() bool {
	return n.recursiveHasFragment()
}

func (n *Node) FragmentType() string {
	if n.Type != "" {
		return n.Type
	}

	for _, child := range n.Children {
		if child.HasFragment() {
			return child.FragmentType()
		}
	}

	return ""
}

func (n *Node) recursiveHasFragment() bool {
	if len(n.Children) == 0 {
		return false
	}

	if n.Type != "" {
		return true
	}

	for _, child := range n.Children {
		if child.recursiveHasFragment() {
			return true
		}
	}

	return false
}

func PlanExecution(selections []query.Selection, fragmentDefinitions query.FragmentDefinitions) []*Node {
	ret := make([]*Node, 0, len(selections))
	for _, sel := range selections {
		switch s := sel.(type) {
		case *query.Field:
			node := &Node{
				Name:       string(s.Name),
				Directives: s.Directives,
				Children:   make([]*Node, 0, len(s.Selections)),
				Arguments:  s.Arguments,
			}

			for _, child := range s.Selections {
				node.Children = append(node.Children, digExecution(child, fragmentDefinitions))
			}

			ret = append(ret, node)
		}
	}

	return ret
}

func digExecution(selectSet query.Selection, fragmentDefinitions query.FragmentDefinitions) *Node {
	switch s := selectSet.(type) {
	case *query.Field:
		node := &Node{
			Name:       string(s.Name),
			Directives: s.Directives,
			Arguments:  s.Arguments,
			Children:   make([]*Node, 0, len(s.Selections)),
		}
		for _, child := range s.Selections {
			switch c := child.(type) {
			case *query.Field:
				node.Children = append(node.Children, digExecution(c, fragmentDefinitions))
			case *query.InlineFragment:
				node.Type = string(c.TypeCondition)
				node.Children = append(node.Children, digExecution(c, fragmentDefinitions))
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
		}
		for _, child := range s.Selections {
			switch c := child.(type) {
			case *query.Field:
				node.Children = append(node.Children, digExecution(c, fragmentDefinitions))
			case *query.InlineFragment:
				node.Children = append(node.Children, digExecution(c, fragmentDefinitions))
			case *query.FragmentSpread:
				node.Children = append(node.Children, digExecution(c, fragmentDefinitions))
			}
		}
		return node
	case *query.FragmentSpread:
		fragment := fragmentDefinitions.GetFragment(s.Name)
		node := &Node{
			Directives: s.Directives,
			Type:       string(fragment.BasedTypeName),
			Children:   make([]*Node, 0, len(fragment.Selections)),
		}
		for _, child := range fragment.Selections {
			switch c := child.(type) {
			case *query.Field:
				node.Children = append(node.Children, digExecution(c, fragmentDefinitions))
			case *query.InlineFragment:
				node.Children = append(node.Children, digExecution(c, fragmentDefinitions))
			case *query.FragmentSpread:
				node.Children = append(node.Children, digExecution(c, fragmentDefinitions))
			}
		}
		return node
	}

	return nil
}
