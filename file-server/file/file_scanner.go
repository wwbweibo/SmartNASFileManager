package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fileserver/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ScanOptions struct {
	Path       []string
	RegexPath  []*regexp.Regexp
	Extensions []string
}

func (opts ScanOptions) OptionPlainPath(path ...string) ScanOptions {
	opts.Path = append(opts.Path, path...)
	return opts
}

func (opts ScanOptions) OptionExtensions(ext ...string) ScanOptions {
	opts.Extensions = append(opts.Extensions, ext...)
	return opts
}

func (opts ScanOptions) OptionRegexPath(regex ...string) ScanOptions {
	for _, r := range regex {
		opts.RegexPath = append(opts.RegexPath, regexp.MustCompile(r))
	}
	return opts
}

func (opts ScanOptions) fileInPath(file string) bool {
	for _, p := range opts.Path {
		if strings.HasPrefix(file, p) {
			return true
		}
	}
	return false
}

func (opts ScanOptions) fileInRegexPath(file string) bool {
	for _, r := range opts.RegexPath {
		if r.MatchString(file) {
			return true
		}
	}
	return false
}

func (opts ScanOptions) fileInExtensions(file string) bool {
	for _, ext := range opts.Extensions {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	return false
}

func ScanAndUpdateFiles(ctx context.Context, path string, option ScanOptions, c chan string) {
	files := utils.WalkDir(path)
	log.Default().Printf("found %d files", len(files))
	for _, file := range files {
		// check if file in option path
		if option.fileInPath(file) {
			c <- file
			continue
		}
		// check if file in option regex path
		if option.fileInRegexPath(file) {
			c <- file
			continue
		}
		// check if file in option extensions
		if option.fileInExtensions(file) {
			c <- file
			continue
		}
	}
}

func StartFileScanner(ctx context.Context, c chan string, repo IFileRepository) {
	for {
		select {
		case path := <-c:
			singleFileHandler(ctx, path, repo)
		case <-ctx.Done():
			return
		}
	}
}

func singleFileHandler(ctx context.Context, file string, repo IFileRepository) {
	log.Default().Printf("handling file %s", file)
	// TODO: invoke third party server to determine the file type
	data := bytes.NewBuffer([]byte(fmt.Sprintf(`{"path": "%s"}`, file)))
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8081/api/v1/file/interfer", data)
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Default().Printf("error getting file type: %v", err)
		return
	}
	defer resp.Body.Close()
	type response struct {
		Type        string `json:"type"`
		Group       string `json:"group"`
		Description string `json:"description"`
	}
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Default().Printf("error reading response body: %v", err)
		return
	}
	var r response
	json.Unmarshal(bts, &r)
	// insert into database
	_file := NewFile(file)
	_file.SetFileType(r.Type, r.Group, r.Description)
	err = repo.CreateOrUpdateFile(ctx, _file)
	if err != nil {
		log.Default().Printf("error inserting file %s: %v", file, err)
	}
}
