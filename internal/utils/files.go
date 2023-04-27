package utils

import (
	"botgpt/internal/config"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetUploadDir() string {
	c := config.GetConfig()

	if runtime.GOOS == "linux" {

		return c.GetString("upload.root_path")
	}
	rootPath := fmt.Sprintf("%s%s", config.GetProjectDir(), c.GetString("upload.root_path"))

	if err := os.Mkdir(rootPath, 0755); os.IsExist(err) {
		// triggers if dir already exists
		return rootPath
	}
	_ = os.Mkdir(rootPath, 0777)

	return rootPath
}

func GetNewFilePath(ext string) (string, string) {
	filePath := GetUploadDir()
	newFileName := strings.ReplaceAll(uuid.New().String(), "-", "")
	return newFileName + ext, filePath + "/" + newFileName + ext
}

func MakeDir(path string) {

	path, err := filepath.Abs(path)
	if err != nil {

	}
	if err := os.Mkdir(path, 0755); os.IsExist(err) {
		// triggers if dir already exists
	}

}
