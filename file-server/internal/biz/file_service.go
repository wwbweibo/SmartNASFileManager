package biz

import (
	"context"
	viewmodel "fileserver/internal/controllers/view_model"
	"fileserver/internal/domain/file"
	"fmt"
	"strings"
)

type FilerService struct {
	repo    file.IFileRepository
	dirNode *dirNode
}

func NewFilerService(repo file.IFileRepository) *FilerService {
	return &FilerService{repo: repo}
}

func (f *FilerService) ListDirectoryTree(ctx context.Context) (viewmodel.ListDirectoryResponse, error) {
	dirs, err := f.repo.ListDirectory(ctx)
	if err != nil {
		return viewmodel.ListDirectoryResponse{}, err
	}
	fmt.Printf("ListDirectoryTree: %v\n", dirs)
	// 构建目录树，使用dirs去构建一个前缀树，输出的结构就是最终的目录树
	var root = dirNode{
		Name:     "/",
		Path:     "/",
		Children: []*dirNode{},
	}
	for _, dir := range dirs {
		root.add(dir)
	}
	f.dirNode = &root
	var resp = f.dirTree2Resp(&root)
	return resp, nil
}

func (f *FilerService) ListFiles(ctx context.Context, path string) ([]viewmodel.ListFileResponseItem, error) {
	if f.dirNode == nil {
		f.ListDirectoryTree(ctx)
	}
	dir := f.dirNode.search(path)
	if dir == nil {
		return nil, nil
	}
	fmt.Printf("ListFiles: %s\n", path)
	files, err := f.repo.ListFileByDirectory(ctx, path)
	if err != nil {
		return nil, err
	}
	var res []viewmodel.ListFileResponseItem
	for _, subDir := range dir.Children {
		res = append(res, viewmodel.ListFileResponseItem{
			Path:  subDir.Path,
			Name:  subDir.Name,
			Type:  "dir",
			Group: "dir",
			Size:  0,
		})
	}
	for _, file := range files {
		res = append(res, viewmodel.ListFileResponseItem{
			Path:  file.Path,
			Name:  file.Name,
			Type:  file.Type,
			Group: file.Group,
			Size:  file.Size,
		})
	}
	return res, nil
}

func (f *FilerService) dirTree2Resp(node *dirNode) viewmodel.ListDirectoryResponse {
	var resp = viewmodel.ListDirectoryResponse{}
	resp.Path = node.Path
	resp.Name = node.Name
	for _, child := range node.Children {
		resp.Children = append(resp.Children, f.dirTree2Resp(child))
	}
	return resp
}

type dirNode struct {
	Name     string
	Path     string
	Children []*dirNode
}

func (node *dirNode) add(path string) {
	section := strings.Split(path, "/")[1:]
	if len(section) == 0 {
		return
	}
	node.append(section[:])
}

func (node *dirNode) append(section []string) {
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
	node.Children = append(node.Children, &dirNode{
		Name: section[0],
		Path: node.Path + section[0] + "/",
	})
}

func (node *dirNode) search(path string) *dirNode {
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

func (node *dirNode) searchNode(section []string) *dirNode {
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
