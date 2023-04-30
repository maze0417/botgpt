package elevenlabs

import (
	"botgpt/internal/config"
	"botgpt/internal/utils"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"github.com/taigrr/elevenlabs/client"
	"github.com/taigrr/elevenlabs/client/types"
)

func Test_TTS(t *testing.T) {
	config.InitTest("botgpt")

	c := config.GetConfig()
	token := c.GetString("elevenlabs.xi_api_key")
	if len(token) == 0 {
		t.Error("can not get token")
		return
	}
	ctx := context.Background()
	// load in an API key to create a client
	client := client.New(token)
	// fetch a list of voice IDs from elevenlabs
	ids, err := client.GetVoiceIDs(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", ids)
	// prepare a pipe for streaming audio directly to beep
	pipeReader, pipeWriter := io.Pipe()
	//reader := bufio.NewReader(os.Stdin)
	//text, _ := reader.ReadString('\n')
	text := "The only true wisdom is in knowing you know nothing.\" - Socrates \n- \"The greatest glory in living lies not in never falling, but in rising every time we fall.\" - Nelson Mandela \n- \"The only way to do great work is to love what you do.\" - Steve Jobs \n- \"In three words I can sum up everything I've learned about life: it goes on.\" - Robert Frost \n- \"Not everything that is faced can be changed, but nothing can be changed until it is faced.\" - James Baldwin \n\nI hope these quotes provide some inspiration and insight for you!"
	go func() {
		// stream audio from elevenlabs using the first voice we found
		err = client.TTSStream(ctx, pipeWriter, text, ids[0], types.SynthesisOptions{Stability: 0.75, SimilarityBoost: 0.75})
		if err != nil {
			panic(err)
		}
		pipeWriter.Close()
	}()
	// decode and prepare the streaming mp3 as it comes through
	streamer, format, err := mp3.Decode(pipeReader)
	if err != nil {
		t.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	// play the audio
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func Test_UploadVoice(t *testing.T) {
	config.InitTest("botgpt")

	c := config.GetConfig()
	token := c.GetString("elevenlabs.xi_api_key")
	if len(token) == 0 {
		t.Error("can not get token")
		return
	}
	ctx := context.Background()
	// load in an API key to create a client
	client := client.New(token)
	// fetch a list of voice IDs from elevenlabs
	voiceFile := utils.GetUploadDir() + "maze.m4a"

	var files []*os.File
	file, err := os.Open(voiceFile)
	if err != nil {
		t.Error("Error opening file:", err)
		return
	}
	files = append(files, file)
	err = client.CreateVoice(ctx, "maze", "", []string{"maze"}, files)
	if err != nil {
		t.Error("Error upload voice file:", err)
		return
	}

	fmt.Println("upload ok")

}
