package file

import (
	"context"
	"fileserver/internal/adapters/dl"
	"fileserver/utils"
	"log"
	"regexp"
	"strings"
	"sync"
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
	_file := NewFile(file)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		result, err := dl.NewClient(dl.Config{
			Scheme: "http",
			Host:   "localhost",
			Port:   8081,
		}).Understanding(ctx, dl.UnderstandingRequest{
			Path: file,
		})
		if err != nil {
			log.Default().Printf("error getting file type: %v", err)
			return
		}
		_file.SetFileTypeFromUnderstanding(result)
	}()
	go func() {
		defer wg.Done()
		_file.Checksum = utils.Sha256(file)
	}()
	// insert into database
	wg.Wait()
	err := repo.CreateOrUpdateFile(ctx, _file)
	if err != nil {
		log.Default().Printf("error inserting file %s: %v", file, err)
	}
}
