package file

import (
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
	node.Children = append(node.Children, &DirNode{
		Name: section[0],
		Path: node.Path + section[0] + "/",
	})
}

func (node *DirNode) Search(path string) *DirNode {
	return node.searchNode(path)
}

func (node *DirNode) searchNode(path string) *DirNode {
	// 做一个 dfs 去搜索所有节点找到对应的节点
	if node.Path == path {
		return node
	}
	for _, n := range node.Children {
		node := n.searchNode(path)
		if node != nil {
			return node
		}
	}
	return nil
}
