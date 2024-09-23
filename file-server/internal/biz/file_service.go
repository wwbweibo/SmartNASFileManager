package biz

import (
	"context"
	viewmodel "fileserver/internal/controllers/view_model"
	"fileserver/internal/domain/file"
	"fmt"
	"strings"
)

type FilerService struct {
	repo file.IFileRepository
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
	var resp = f.dirTree2Resp(&root)
	return resp, nil
}

func (f *FilerService) ListFiles(ctx context.Context, path string) ([]viewmodel.ListFileResponseItem, error) {
	fmt.Printf("ListFiles: %s\n", path)
	files, err := f.repo.ListFileByDirectory(ctx, path)
	if err != nil {
		return nil, err
	}
	var res []viewmodel.ListFileResponseItem
	for _, file := range files {
		res = append(res, viewmodel.ListFileResponseItem{
			Path: file.Path,
			Name: file.Name,
			Type: file.Type,
			Size: file.Size,
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
	fmt.Printf("section: ||%s||\n", strings.Join(section, "--"))
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
