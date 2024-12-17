package viewmodel

type PageRequest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

func (req PageRequest) GetOffset() int {
	return (req.Page - 1) * req.Size
}
