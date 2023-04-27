package aws

import (
	"botgpt/internal/enum"
	"botgpt/internal/interfaces"
	"github.com/pemistahl/lingua-go"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

var pollyClient *polly.Polly
var once sync.Once

var LangMap = map[string]string{
	enum.JaJP:  "Mizuki",
	enum.EnUS:  "Joanna",
	enum.CmnCN: "Zhiyu",
}

type Polly struct {
}

func NewPolly() interfaces.ITextToSpeech {
	return &Polly{}
}

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
func (p Polly) GetLangFromText(text string) string {
	detector := lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		Build()

	language, exists := detector.DetectLanguageOf(text)

	if !exists {
		return ""
	}

	lang, ok := enum.LangMap[language.String()]
	if !ok {
		return language.String()
	}
	return lang
}

func (p Polly) TextToSpeech(text string, outputFile string, outputFormat string, lang string) error {
	return textToSpeech(text, outputFile, outputFormat, lang)
}

func textToSpeech(text string, outputFile string, outputFormat string, lang string) error {
	client := getPollyClient()

	actor, ok := LangMap[lang]

	if !ok {
		//err := fmt.Sprintf("No speech lang info for %s \n", lang)
		//log.Printf(err)
		//return errors.New(err)
		lang = enum.CmnCN
		actor = LangMap[enum.CmnCN]
	}

	log.Printf("Send to Polly Actor is %s for %s \n", actor, lang)
	// Synthesize speech using AWS Polly
	input := &polly.SynthesizeSpeechInput{
		Engine:          nil,
		LanguageCode:    aws.String(lang),
		LexiconNames:    nil,
		OutputFormat:    aws.String(outputFormat),
		SampleRate:      nil,
		SpeechMarkTypes: nil,
		Text:            aws.String(text),
		TextType:        nil,
		VoiceId:         aws.String(actor), // Replace "Joanna" with the desired voice
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
