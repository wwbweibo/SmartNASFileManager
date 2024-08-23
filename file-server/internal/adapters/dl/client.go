package dl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	_understandingPath = "/api/v1/file/understanding"
)

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Understanding(ctx context.Context, req UnderstandingRequest) (r UnderstandingResult, err error) {
	reqBts, _ := json.Marshal(req)
	data := bytes.NewBuffer(reqBts)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s://%s:%d%s", c.config.Scheme, c.config.Host, c.config.Port, _understandingPath), data)
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Default().Printf("error getting file type: %v", err)
		return
	}
	defer resp.Body.Close()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Default().Printf("error reading response body: %v", err)
		return
	}
	json.Unmarshal(bts, &r)
	return
}
