package viewmodel

type ListFileRequest struct {
	Path string `json:"path" form:"path"`
}

type ListFileResponseItem struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type ListDirectoryResponse struct {
	Name     string                  `json:"name"`
	Path     string                  `json:"path"`
	Children []ListDirectoryResponse `json:"children"`
}
