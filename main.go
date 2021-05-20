package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"dbhushan9/vaccine-alerts/bot"
	"dbhushan9/vaccine-alerts/types"
	arrayUtils "dbhushan9/vaccine-alerts/util"

	"github.com/jasonlvhit/gocron"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	indianLocale = "Asia/Kolkata"
)

//go:embed "templates/telegram-notification.md"
var telegramMsgTemplate string

//go:embed "templates/email-template.html"
var emailTemplate string

func main() {
	// setLoggingConfig()
	log.Print("Starting Vaccine alert worker")
	// load .env file
	vaccineDate := arrayUtils.GetNextLocaleDay(indianLocale)
	ageSlots := [2]int{18, 45}
	districtId := 363
	blockName := "Haveli"
	vaccine := "any"
	feeType := "any"

	channelID := os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT")
	debugChannelID := os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT_DEBUG")

	var emails = []string{"dbhushan912@gmail.com", "avi6387@gmail.com", "deshmukh.kalyani81@gmail.com", "rupadeshmukh26@gmail.com"}

	//notificationOptions
	//vaccineSearchOptions
	gocron.Every(5).Minute().Do(func() {
		checkForVaccineCenters(vaccineDate, districtId, ageSlots, blockName, vaccine, feeType, channelID, debugChannelID, emails)
	})
	<-gocron.Start()
	// checkForVaccineCenters(vaccineDate, districtId)
}

func sendMail(msg string, name string, email string) {
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

func checkForVaccineCenters(vaccineDate string, districtId int, ageSlots [2]int, blockName string, vaccine string, feeType string, channelID string, debugChannelID string, emails []string) {
	log.Printf("Checking for Vaccine Centers for date: %v districtId:%d", vaccineDate, districtId)
	client := &http.Client{}
	url := fmt.Sprintf("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?date=%v&district_id=%d", vaccineDate, districtId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("Could not create Http request")
		log.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "hi_IN")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to %v", url)
		log.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error reading request body")
		log.Print(err.Error())
	}

	var response types.ReponseType
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Print("Error unmarshalling request body")
		log.Print(err.Error())
	}
	log.Printf("Total centers are %v", len(response.Centers))

	_, cd18 := types.ProcessCentersPresent(response.Centers, ageSlots[0], vaccine, feeType, vaccineDate, blockName)
	_, cd45 := types.ProcessCentersPresent(response.Centers, ageSlots[1], vaccine, feeType, vaccineDate, blockName)

	centerLogs := types.LogDetails{Age18: *cd18, Age45: *cd45}
	logJSON, _ := json.MarshalIndent(centerLogs, "", " ")

	if len(*cd18) > 0 || len(*cd45) > 0 {
		msg := renderTemplate(centerLogs, vaccineDate, telegramMsgTemplate)
		bot.SendTelegramMessage(msg, true, channelID)

		msg = fmt.Sprintf("%v - Found Available centers for %v", time.Now().Format("01-02-2006 15:04:05"), vaccineDate)
		bot.SendTelegramMessage(msg, false, debugChannelID)

		// emailMsg := renderTemplate(centerLogs, vaccineDate, emailTemplate)
		// for _, mail := range emails {
		// 	sendMail(emailMsg, "User", mail)
		// }

		log.Printf("Centers Available")
		_ = ioutil.WriteFile(fmt.Sprintf("%v-centers-%v.json", vaccineDate, time.Now().Format("01-02-2006-15-04-05")), logJSON, 0644)
	} else {
		msg := fmt.Sprintf("%v - No Available centers for %v", time.Now().Format("01-02-2006 15:04:05"), vaccineDate)
		bot.SendTelegramMessage(msg, false, debugChannelID)
		log.Print("Centers Not Availabe")
	}
}

func renderTemplate(centers types.LogDetails, vaccineDate string, templateString string) string {
	student := struct {
		Date           string
		TotalCenters45 int
		TotalCenters18 int
		CenterDetails  types.LogDetails
	}{
		Date:           vaccineDate,
		TotalCenters45: len(centers.Age45),
		TotalCenters18: len(centers.Age18),
		CenterDetails:  centers,
	}

	parsedTemplate, err := template.New("msg-template").Parse(templateString)
	if err != nil {
		panic(err)
	}

	var templateBytes strings.Builder
	_ = parsedTemplate.Execute(&templateBytes, student)
	return templateBytes.String()
}

func setLoggingConfig() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	log.SetOutput(file)
	log.Print("Logging to a file in Go!")

}
