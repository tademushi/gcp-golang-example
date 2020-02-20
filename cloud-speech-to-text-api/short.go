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
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	filename := "short.ogg"

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_OGG_OPUS,
			SampleRateHertz: 48000,
			LanguageCode:    "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
}
