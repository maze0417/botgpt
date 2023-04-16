package aws

import (
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

var pollyClient *polly.Polly
var once sync.Once

func getPollyClient() *polly.Polly {
	once.Do(func() {
		// Create a new session with the default AWS configuration
		sess, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			panic(err)
		}

		// Create a Polly client
		pollyClient = polly.New(sess)
	})

	return pollyClient
}
func SynthesizeSpeech(text string, outputFile string, outputFormat string) error {
	client := getPollyClient()

	// Synthesize speech using AWS Polly
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String(outputFormat),
		Text:         aws.String(text),
		VoiceId:      aws.String("Joanna"), // Replace "Joanna" with the desired voice
	}

	output, err := client.SynthesizeSpeech(input)
	if err != nil {
		return err
	}

	// Read the output audio stream
	audioBytes, err := io.ReadAll(output.AudioStream)
	if err != nil {
		return err
	}

	// Save the synthesized speech to a file
	err = os.WriteFile(outputFile, audioBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
