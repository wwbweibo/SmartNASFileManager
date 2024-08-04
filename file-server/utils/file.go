package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
)

func WalkDir(dir string) []string {
	var _files []string
	// check if the dir exist on the filesystem
	log.Default().Printf("walking dir %s", dir)
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Default().Printf("dir %s does not exist", dir)
			return []string{}
		} else {
			log.Default().Fatalf("error reading dir %s: %v", dir, err)
		}
	}
	// check if the dir is a dir
	if info.IsDir() {
		// list dir
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Default().Printf("error reading dir %s: %v", dir, err)
			return []string{}
		}
		for _, f := range files {
			if f.IsDir() {
				// recursively walk the dir
				_files = append(_files, WalkDir(path.Join(dir, f.Name()))...)
			} else {
				// append the file to the list
				_files = append(_files, path.Join(dir, f.Name()))
			}
		}
	} else {
		_files = append(_files, dir)
	}
	return _files
}

func GetDirectory(f string) string {
	return f[:len(f)-len(GetFileName(f))]
}

func GetFileName(f string) string {
	for i := len(f) - 1; i >= 0; i-- {
		if f[i] == '/' {
			return f[i+1:]
		}
	}
	return f
}

func GetExtension(f string) string {
	for i := len(f) - 1; i >= 0; i-- {
		if f[i] == '.' {
			return f[i:]
		}
	}
	return ""
}

func Sha256(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Default().Fatalf("error opening file %s: %v", path, err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Default().Fatalf("error calculating SHA256 for file %s: %v", path, err)
	}

	return hex.EncodeToString(hash.Sum(nil))
}
