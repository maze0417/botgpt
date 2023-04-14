package ai

import (
	"fmt"
	"testing"
)

func TestTextToSpeech(t *testing.T) {

	text := "Hello, I am using AWS Polly with Golang."
	outputFile := "testts.mp3"

	err := synthesizeSpeech(text, outputFile)
	if err != nil {

		t.Errorf("Failed to synthesize speech: %v\n", err)

	}

	fmt.Println("Synthesized speech saved to", outputFile)

}
