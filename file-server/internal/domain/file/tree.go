package file

import (
	"fmt"
	"strings"
)

var Root = &DirNode{
	Name:     "/",
	Path:     "/",
	Children: []*DirNode{},
}

type DirNode struct {
	Name     string
	Path     string
	Children []*DirNode
}

func (node *DirNode) Add(path string) {
	section := strings.Split(path, "/")[1:]
	if len(section) == 0 {
		return
	}
	node.append(section[:])
}

func (node *DirNode) append(section []string) {
	if len(section) == 0 {
		return
	}
	if len(section) == 1 && section[0] == "" {
		return
	}
	for _, n := range node.Children {
		if n.Name == section[0] {
			n.append(section[1:])
			return
		}
	}
	_node := &DirNode{
		Name: section[0],
		Path: node.Path + section[0] + "/",
	}
	_node.append(section[1:])
	node.Children = append(node.Children, _node)
}

func (node *DirNode) Search(path string) *DirNode {
	return node.searchNode(path)
}

func (node *DirNode) searchNode(path string) *DirNode {
	// 做一个 dfs 去搜索所有节点找到对应的节点
	nodes := []*DirNode{node}
	for len(nodes) > 0 {
		n := nodes[0]
		nodes = nodes[1:]
		fmt.Printf("searching %s\n", n.Path)
		fmt.Printf("'%s' == '%s' = %v\n", n.Path, path, n.Path == path)
		if n.Path == path {
			fmt.Printf("found %s\n", n.Path)
			return n
		}
		if len(n.Children) > 0 {
			nodes = append(nodes, n.Children...)
		}
	}
	return nil
}
