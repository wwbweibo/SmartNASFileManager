package file

import (
	"context"

	"gorm.io/gorm"
)

type IFileRepository interface {
	ListFileByDirectory(ctx context.Context, directory string) ([]File, error)
	CreateOrUpdateFile(ctx context.Context, file File) (err error)
}

type FileQuery struct {
	Path      string
	Directory string
	Extension string
	FileType  string
}

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	db.AutoMigrate(&File{})
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) ListFileByDirectory(ctx context.Context, directory string) ([]File, error) {
	var files []File
	err := r.db.Where("directory = ?", directory).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) CreateOrUpdateFile(ctx context.Context, file File) (err error) {
	return r.db.Exec(`insert into files 
	(path, directory, extension, type, "group", description, tags ) values 
	(?,?,?,?,?,?,?) on conflict(path) do update set directory = ?, extension = ?, type = ?, "group"=?, description=? , tags = ?;`,
		file.Path, file.Directory, file.Extension, file.Type, file.Group, file.Description, file.Tags,
		file.Directory, file.Extension, file.Type, file.Group, file.Description, file.Tags).Error
}
