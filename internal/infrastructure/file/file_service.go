package file

import (
	"fmt"
)

type FileService struct {
	Dir     string
	BaseUrl string
}

func NewFileService(dir string) *FileService {
	return &FileService{Dir: dir}
}

// GetFilePath 获取文件路径
func (s *FileService) GetFilePath(fileName string) string {
	return fmt.Sprintf("%s_%s", s.Dir, fileName)
}

// GetFileURL 获取文件URL
func (s *FileService) GetFileURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", s.BaseUrl, filePath)
}
