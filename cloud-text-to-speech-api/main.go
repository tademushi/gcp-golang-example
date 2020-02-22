// オリジナルのソース
// https://cloud.google.com/text-to-speech/docs/quickstart-client-libraries?hl=ja

package main

import (
	"context"
	"io/ioutil"
	"log"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {
	ctx := context.Background()
	text := "Its a sunny day"
	languageCode := "en-US"
	voiceGender := texttospeechpb.SsmlVoiceGender_NEUTRAL
	audioEncoding := texttospeechpb.AudioEncoding_OGG_OPUS
	audioData, err := createSpeech(ctx, text, languageCode, voiceGender, audioEncoding)
	if err != nil {
		log.Fatal(err)
	}

	filename := "output.ogg"
	err = writeFile(filename, audioData)
	if err != nil {
		log.Fatal(err)
	}
}

// createSpeech は文字列から音声データを生成します
// voiceGender https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/texttospeech/v1?tab=doc#SsmlVoiceGender
// audioEncoding https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/texttospeech/v1?tab=doc#AudioEncoding
func createSpeech(context context.Context, text string, languageCode string, voiceGender texttospeechpb.SsmlVoiceGender, audioEncoding texttospeechpb.AudioEncoding) ([]byte, error) {
	client, err := texttospeech.NewClient(context)
	if err != nil {
		return nil, err
	}

	request := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: languageCode,
			SsmlGender:   voiceGender,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: audioEncoding,
		},
	}

	response, err := client.SynthesizeSpeech(context, &request)
	if err != nil {
		return nil, err
	}

	return response.AudioContent, nil
}

func writeFile(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, 0644)
	return err
}
