package file

import (
	"encoding/json"
	"fileserver/internal/adapters/dl"
	"fileserver/utils"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Path        string `gorm:"uniqueindex"`
	Name        string
	Directory   string
	Extension   string
	Type        string
	Group       string
	Description string
	Tags        string
	Caption     string
	Checksum    string
	Size        int64
}

func NewFile(path string) File {
	// calc sha256 hash checksum
	return File{
		Path:      path,
		Name:      utils.GetFileName(path),
		Directory: utils.GetDirectory(path),
		Extension: utils.GetExtension(path),
	}
}

func (f *File) SetFileType(ttype, group, description string) {
	f.Type = ttype
	f.Group = group
	f.Description = description
}

func (f *File) SetFileTypeFromUnderstanding(understanding dl.UnderstandingResult) {
	f.Type = understanding.Label
	f.Group = understanding.Group
	f.Description = understanding.Description
	f.setupFileExtensionInfo(understanding)
}

func (f *File) setupFileExtensionInfo(understanding dl.UnderstandingResult) {
	if understanding.Extension == nil {
		return
	}
	if understanding.Group == "image" {
		bts, _ := json.Marshal(understanding.Extension)
		ext := imageUnderstandingExtension{}
		json.Unmarshal(bts, &ext)
		f.Caption = ext.Caption
		if len(ext.Labels) > 0 {
			var tags []string
			linq.From(ext.Labels).SelectT(func(lable imageUnderstandingExtensionLabel) string {
				return lable.Label
			}).ToSlice(&tags)
			f.Tags = strings.Join(tags, ",")
		}
	}
}

func (f *File) CalcSha256() {
	f.Checksum = utils.Sha256(f.Path)
}

type imageUnderstandingExtension struct {
	Caption string                             `json:"caption"`
	Labels  []imageUnderstandingExtensionLabel `json:"labels"`
}

type imageUnderstandingExtensionLabel struct {
	Label      string `json:"label"`
	Confidence string `json:"confidence"`
}
