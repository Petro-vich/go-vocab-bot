package translation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-vocab-bot/internal/domain"
	"net/http"
	"os"
)

type YandexTranslator struct {
	apiKey   string
	folderID string
}

func NewTranslator() domain.Translator {
	return &YandexTranslator{
		apiKey:   os.Getenv("YANDEX_API_KEY"),
		folderID: os.Getenv("YANDEX_FOLDER_ID"),
	}
}

func (yt *YandexTranslator) Translate(texts string, targetLang string) (string, error) {
	requestBody := map[string]interface{}{
		"folderId":           yt.folderID,
		"texts":              []string{texts},
		"targetLanguageCode": targetLang,
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(
		"POST",
		"https://translate.api.cloud.yandex.net/translate/v2/translate",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Authorization", "Api-Key "+yt.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var result struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}

	var translations []string
	for _, t := range result.Translations {
		translations = append(translations, t.Text)
	}

	return result.Translations[0].Text, nil
}
