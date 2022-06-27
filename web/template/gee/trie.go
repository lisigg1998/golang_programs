package gee

import "strings"

type node struct {
	pattern  string  // the URL path to match. only leaf nodes set this field.
	part     string  // the part to match in this layer
	children []*node // pointers to children nodes in next layer。 TODO: how about making it a hashset, or rb-tree?
	isWild   bool    // whether it is wildcard ":" or "*"
}

// return the first matched child
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// return all matched children
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.isWild || child.part == part {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert to trie
// TODO: this will insert everything to wildcard node. need to optimize later.
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search the given parts
// TODO: also need to correct logic and optimize. the key point is, distinguish "routing rule" and "actual route path"
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern != "" {
			return n
		}
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
