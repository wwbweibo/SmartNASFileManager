package dl

type UnderstandingRequest struct {
	Path string `json:"path"`
}

type UnderstandingResult struct {
	Label       string `json:"label"`
	Group       string `json:"group"`
	Description string `json:"description"`
	Extension   any    `json:"extension"`
}
