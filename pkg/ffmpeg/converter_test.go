package ffmpeg

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCanConvert(t *testing.T) {

	_, sourceFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Errorf("Error getting source file path")

		return
	}
	projectFolder := filepath.Dir(filepath.Dir(sourceFilePath))

	inputOggName := fmt.Sprintf("%s/%s-%d.ogg", projectFolder, "testgroup", 123)
	outputMp3Name := fmt.Sprintf("%s/%s-%d.mp3", projectFolder, "testgroup", 123)

	err := ConvertOggToMp3(inputOggName, outputMp3Name)
	if err != nil {
		t.Error(err)

		return
	}

}
