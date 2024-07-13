package file

import (
	"fileserver/utils"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Path        string `gorm:"uniqueindex"`
	Directory   string
	Extension   string
	Type        string
	Group       string
	Description string
	Tags        string
}

func NewFile(path string) File {
	return File{
		Path:      path,
		Directory: utils.GetDirectory(path),
		Extension: utils.GetExtension(path),
	}
}

func (f *File) SetFileType(ttype, group, description string) {
	f.Type = ttype
	f.Group = group
	f.Description = description
}
