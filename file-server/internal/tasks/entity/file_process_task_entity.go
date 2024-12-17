package entity

// FileProcessTask 文件处理任务
type FileProcessTask struct {
	File string
}

// GetTaskName 获取任务名称
func (t *FileProcessTask) GetTaskName() string {
	return "FileProcessTask"
}
