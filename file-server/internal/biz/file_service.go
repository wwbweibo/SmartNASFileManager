package biz

import (
	"context"
	viewmodel "fileserver/internal/controllers/view_model"
	"fileserver/internal/domain/file"
	"fmt"
)

type FilerService struct {
	repo file.IFileRepository
}

func NewFilerService(repo file.IFileRepository) *FilerService {
	return &FilerService{repo: repo}
}

func (f *FilerService) ListDirectoryTree(ctx context.Context) (viewmodel.ListDirectoryResponse, error) {
	var resp = f.dirTree2Resp(file.Root)
	return resp, nil
}

func (f *FilerService) ListFiles(ctx context.Context, path string) ([]viewmodel.ListFileResponseItem, error) {
	dir := file.Root.Search(path)
	if dir == nil {
		fmt.Printf("list file in '%s' error, could not find directory\n", path)
		return nil, nil
	}
	fmt.Printf("ListFiles: %s\n", path)
	files, err := f.repo.ListFileByDirectory(ctx, path)
	if err != nil {
		fmt.Printf("list file in %s error: %v\n", path, err.Error())
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

func (f *FilerService) dirTree2Resp(node *file.DirNode) viewmodel.ListDirectoryResponse {
	var resp = viewmodel.ListDirectoryResponse{}
	resp.Path = node.Path
	resp.Name = node.Name
	for _, child := range node.Children {
		resp.Children = append(resp.Children, f.dirTree2Resp(child))
	}
	return resp
}
