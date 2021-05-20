package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(msg string, name string, email string) {
	from := mail.NewEmail("Vaccine Alerts", os.Getenv("SENDER_EMAIL"))
	subject := "Found Vaccine matching your criteria"
	to := mail.NewEmail(name, email)
	message := mail.NewSingleEmail(from, subject, to, "", msg)
	client := sendgrid.NewSendClient(os.Getenv("APIKEY_SENDGRID"))
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Sent Email to %v", email)
	}
}

func SendTelegramNotification(msg string, isParseMode bool, channelID string) {
	log.Print("Sending Telegram Message")
	apiKey := os.Getenv("APIKEY_TELEGRAM_BOT")
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", apiKey)

	requestBody := createTelegramRequestBody(msg, isParseMode, channelID)
	_, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		log.Print(err.Error())
	}
	//TODO: check response status
	log.Print("Telegram message sent to channel")
}

func createTelegramRequestBody(msg string, isParseMode bool, channelID string) *bytes.Buffer {
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
