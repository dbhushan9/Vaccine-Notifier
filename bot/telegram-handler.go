package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func SendTelegramMessage(msg string, isParseMode bool, channelID string) {
	log.Print("Sending Telegram Message")
	apiKey := os.Getenv("APIKEY_TELEGRAM_BOT")
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", apiKey)

	requestBody := createRequestBody(msg, isParseMode, channelID)
	_, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		log.Print(err.Error())
	}
	//TODO: check response status
	log.Print("Telegram message sent to channel")
}

func createRequestBody(msg string, isParseMode bool, channelID string) *bytes.Buffer {
	var postBody []byte
	if isParseMode {
		postBody, _ = json.Marshal(map[string]string{
			"chat_id":    channelID,
			"text":       msg,
			"parse_mode": "MarkdownV2",
		})
	} else {
		postBody, _ = json.Marshal(map[string]string{
			"chat_id": channelID,
			"text":    msg,
		})
	}

	return bytes.NewBuffer(postBody)
}
