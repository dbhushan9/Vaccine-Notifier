package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(msg string, name string, email string) {
	log.Info("sending email")
	from := mail.NewEmail("Vaccine Alerts", os.Getenv("SENDER_EMAIL"))
	subject := "Found Vaccine matching your criteria"
	to := mail.NewEmail(name, email)
	message := mail.NewSingleEmail(from, subject, to, "", msg)
	client := sendgrid.NewSendClient(os.Getenv("APIKEY_SENDGRID"))
	_, err := client.Send(message)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "email": email}).Error("error sending email")
	}
	log.WithFields(log.Fields{"email": email}).Info("email sent")

}

func SendTelegramNotification(msg string, isParseMode bool, channelID string) {
	apiKey := os.Getenv("APIKEY_TELEGRAM_BOT")
	log.WithFields(log.Fields{"channelID": channelID, "bot_token": apiKey, "msg": msg}).Info("sending telegram message")
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", apiKey)

	requestBody := createTelegramRequestBody(msg, isParseMode, channelID)
	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error sending telegram message")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error reading body")
	}

	var data interface{}
	json.Unmarshal(body, data)
	//TODO: check response status
	if resp.StatusCode == 200 {
		log.WithFields(log.Fields{"response_status": resp.StatusCode, "response_body": body}).Info("telegram message sent")
	} else {
		log.WithFields(log.Fields{"response_status": resp.StatusCode, "response_body": body}).Info("error sending telegram message")
	}
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
