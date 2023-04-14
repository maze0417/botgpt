package ai

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

func synthesizeSpeech(text string, outputFile string) error {
	// Create a new session with the default AWS configuration
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	// Create a Polly client
	pollyClient := polly.New(sess)

	// Synthesize speech using AWS Polly
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		Text:         aws.String(text),
		VoiceId:      aws.String("Joanna"), // Replace "Joanna" with the desired voice
	}

	output, err := pollyClient.SynthesizeSpeech(input)
	if err != nil {
		return err
	}

	// Read the output audio stream
	audioBytes, err := ioutil.ReadAll(output.AudioStream)
	if err != nil {
		return err
	}

	// Save the synthesized speech to a file
	err = ioutil.WriteFile(outputFile, audioBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
