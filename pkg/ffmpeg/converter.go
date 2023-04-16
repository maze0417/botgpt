package ffmpeg

import (
	"fmt"
	"os/exec"
)

func ConvertOggToMp3(inputOggName string, outputMp3Name string) error {
	cmd := exec.Command("ffmpeg", "-i", inputOggName, outputMp3Name)
	err := cmd.Run()
	if err != nil {
		fmt.Println("ffmpeg Error converting OGG to MP3:", err)
		return err
	}
	return nil
}

func ConvertMp3ToOgg(inputMo3Name string, outputOggName string) error {
	cmd := exec.Command("ffmpeg", "-i", inputMo3Name, "-c:a", "libopus", outputOggName)
	err := cmd.Run()
	if err != nil {
		fmt.Println("ffmpeg Error converting OGG to MP3:", err)
		return err
	}
	return nil
}
