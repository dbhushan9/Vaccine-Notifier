package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func SendTelegramMessage(msg string, isParseMode bool) {
	log.Print("Sending Telegram Message")
	apiKey := os.Getenv("APIKEY_TELEGRAM_BOT")
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", apiKey)

	requestBody := createRequestBody(msg, isParseMode)

	_, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		log.Print(err.Error())
	}
	//TODO: check response status
	log.Print("Telegram message sent to channel")
}

func createRequestBody(msg string, isParseMode bool) *bytes.Buffer {
	var postBody []byte
	if isParseMode {
		postBody, _ = json.Marshal(map[string]string{
			"chat_id":    os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT"),
			"text":       msg,
			"parse_mode": "MarkdownV2",
		})
	} else {
		postBody, _ = json.Marshal(map[string]string{
			"chat_id": os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT_DEBUG"),
			"text":    msg,
		})
	}

	return bytes.NewBuffer(postBody)
}
