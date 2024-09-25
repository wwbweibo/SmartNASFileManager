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
	fmt.Printf("append %v to %s\n", section, node.Path)
	for _, n := range node.Children {
		if n.Name == section[0] {
			fmt.Printf("dir exist %s ", n)
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
	// 如果 path 以/结束，去掉最后的/
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	section := strings.Split(path, "/")[1:]
	if len(section) == 0 {
		return node
	}
	return node.searchNode(section[:])
}

func (node *DirNode) searchNode(section []string) *DirNode {
	if len(section) == 0 {
		return node
	}
	for _, n := range node.Children {
		if n.Name == section[0] {
			return n.searchNode(section[1:])
		}
	}
	return nil
}
