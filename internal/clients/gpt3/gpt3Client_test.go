package gpt3

import (
	"botgpt/internal/config"
	"fmt"
	"testing"
)

func TestClient_Transcribe(t *testing.T) {
	config.InitTest()

	outputMp3Name := fmt.Sprintf("D:/AgProject/sms_compose/sms_server/ffmpeg/testgroup-123.mp3")
	res, err := Transcribe(outputMp3Name)
	if err != nil {
		t.Error(err)

		return
	}

	fmt.Println(res)

}

func TestM4a(t *testing.T) {
	config.InitTest()

	outputMp3Name := fmt.Sprintf("D:/AgProject/sms_compose/sms_server/ffmpeg/test_noise.m4a")
	res, err := Transcribe(outputMp3Name)
	if err != nil {
		t.Error(err)

		return
	}

	fmt.Println(res)

}
