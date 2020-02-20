package main

import (
	"context"
	"fmt"

	translate "cloud.google.com/go/translate/apiv3"
	"golang.org/x/oauth2/google"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

func getProjectID(ctx context.Context) (string, error) {
	crednetials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return "", err
	}
	return crednetials.ProjectID, nil
}

func createParent(ctx context.Context) (string, error) {
	projectID, err := getProjectID(ctx)
	if err != nil {
		return "", err
	}

	parent := fmt.Sprintf("projects/%s", projectID)

	return parent, nil
}

func translateText(sourceLanguage string, targetLanguage string, text string) (string, error) {
	ctx := context.Background()

	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	parent, err := createParent(ctx)
	if err != nil {
		return "", err
	}

	request := &translatepb.TranslateTextRequest{
		Contents:           []string{text},
		MimeType:           "text/plain",
		SourceLanguageCode: sourceLanguage,
		TargetLanguageCode: targetLanguage,
		Parent:             parent,
	}

	response, err := client.TranslateText(ctx, request)
	if err != nil {
		return "", err
	}

	if len(response.Translations) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %v", text)
	}
	return response.Translations[0].TranslatedText, nil
}

func main() {
	sourceText := "本日は晴天なり"

	translatedText, err := translateText("ja", "en", sourceText)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("sourceText   :", sourceText)
	fmt.Println("translatedText:", translatedText)
}
