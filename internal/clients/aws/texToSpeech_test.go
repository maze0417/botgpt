package aws

import (
	"botgpt/internal/clients/telegram"
	"botgpt/internal/config"
	"botgpt/internal/utils"
	"fmt"
	"github.com/aws/aws-sdk-go/service/polly"
	"testing"
)

func TestTextToSpeech(t *testing.T) {

	config.InitTest("botgpt")

	text := "Hello, I am using AWS Polly with Golang."
	outputFile := utils.GetUploadDir() + "testts." + polly.OutputFormatMp3

	err := SynthesizeSpeech(text, outputFile, polly.OutputFormatMp3)
	if err != nil {

		t.Errorf("Failed to synthesize speech: %v\n", err)
		return
	}

	fmt.Println("Synthesized speech saved to", outputFile)

	err = telegram.SendVoice(outputFile)
	if err != nil {
		t.Error(err)
	}
}
