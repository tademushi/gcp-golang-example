// オリジナルのソース
// https://cloud.google.com/speech-to-text/docs/sync-recognize
//
// 使用可能な音声ファイルの仕様について書かれている
// https://cloud.google.com/speech-to-text/docs/encoding

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {
	filename := "short.ogg"
	audioData, err := readFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	ctx := context.Background()
	encoding := speechpb.RecognitionConfig_OGG_OPUS
	var sampleRateHertz int32 = 48000
	languageCode := "ja-JP"
	textList, err := createText(ctx, audioData, encoding, sampleRateHertz, languageCode)

	for _, item := range textList {
		fmt.Printf("\"%v\" (confidence=%3f)\n", item.text, item.confidence)
	}
}

type speechtotextResult struct {
	text       string
	confidence float32
}

// createText は音声データから文字列を生成します
// encoding https://godoc.org/google.golang.org/genproto/googleapis/cloud/speech/v1#RecognitionConfig_AudioEncoding
func createText(context context.Context, audioData []byte, encoding speechpb.RecognitionConfig_AudioEncoding, sampleRateHertz int32, languageCode string) ([]speechtotextResult, error) {
	client, err := speech.NewClient(context)
	if err != nil {
		return nil, err
	}

	request := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        encoding,
			SampleRateHertz: sampleRateHertz,
			LanguageCode:    languageCode,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioData},
		},
	}

	response, err := client.Recognize(context, request)
	if err != nil {
		return nil, err
	}

	var rtn = make([]speechtotextResult, 0)
	for _, result := range response.Results {
		for _, alt := range result.Alternatives {
			item := speechtotextResult{
				text:       alt.Transcript,
				confidence: alt.Confidence,
			}
			rtn = append(rtn, item)
		}
	}

	return rtn, nil
}

func readFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	return data, err
}
