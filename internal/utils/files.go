package utils

import (
	"botgpt/internal/config"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const FileTimeFormat = "2006-01-02"

func GetFileDir() string {
	c := config.GetConfig()
	rootPath := c.GetString("upload.root_path")
	now := time.Now()
	strTime := strings.Split(now.Format(FileTimeFormat), "-")
	year := strTime[0]
	month := strTime[1]
	day := strTime[2]
	MakeDir(rootPath + year)
	MakeDir(rootPath + year + "/" + month)
	MakeDir(rootPath + year + "/" + month + "/" + day)
	path, err := filepath.Abs(rootPath + year + "/" + month + "/" + day)
	if err != nil {

	}
	if err := os.Mkdir(rootPath, 0755); os.IsExist(err) {
		// triggers if dir already exists
		return path
	} else {
		err = os.Mkdir(rootPath, 0777)
		if err != nil {

		}
	}

	return path
}

func GetNewFilePath(ext string) (string, string) {
	filePath := GetFileDir()
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
