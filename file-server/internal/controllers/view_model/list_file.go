package viewmodel

type ListFileRequest struct {
	Path string `json:"path" form:"path"`
}

type ListFileByGroupRequest struct {
	Group string `json:"group" form:"group"`
	Path  string `json:"path" form:"path"`
	PageRequest
}

type ListFileResponseItem struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Size  int64  `json:"size"`
	Group string `json:"group"`
}

type ListDirectoryResponse struct {
	Name     string                  `json:"name"`
	Path     string                  `json:"path"`
	Children []ListDirectoryResponse `json:"children"`
}
