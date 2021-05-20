package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"dbhushan9/vaccine-alerts/cowin"
	"dbhushan9/vaccine-alerts/notifier"

	"github.com/jasonlvhit/gocron"
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
	vaccineDate := GetNextLocaleDay(indianLocale)
	ageSlots := [2]int{18, 45}
	districtId := 363
	blockName := "Haveli"
	vaccine := "any"
	feeType := "any"
	//ARMY ONLY CENTERS
	excludedCenter := []int{619964, 629727}

	channelID := os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT")
	debugChannelID := os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT_DEBUG")

	var emails = []string{"dbhushan912@gmail.com", "avi6387@gmail.com", "deshmukh.kalyani81@gmail.com", "rupadeshmukh26@gmail.com"}

	//notificationOptions
	//vaccineSearchOptions
	gocron.Every(5).Minute().Do(func() {
		checkForVaccineCenters(vaccineDate, districtId, ageSlots, blockName, vaccine, feeType, channelID, debugChannelID, emails, excludedCenter)
	})
	<-gocron.Start()
	// checkForVaccineCenters(vaccineDate, districtId)
}

func checkForVaccineCenters(vaccineDate string, districtId int, ageSlots [2]int, blockName string, vaccine string, feeType string, channelID string, debugChannelID string, emails []string, excludedCenter []int) {
	log.Printf("Checking for Vaccine Centers for date: %v districtId:%d", vaccineDate, districtId)

	response, err := cowin.QueryCowinAPI(vaccineDate, districtId)
	if err != nil {
		log.Print("Error unmarshalling request body")
		log.Print(err.Error())
	}

	log.Printf("Total centers are %v", len(response.Centers))

	_, cd18 := cowin.ProcessCentersPresent(response.Centers, ageSlots[0], vaccine, feeType, vaccineDate, blockName, excludedCenter)
	_, cd45 := cowin.ProcessCentersPresent(response.Centers, ageSlots[1], vaccine, feeType, vaccineDate, blockName, excludedCenter)

	centerLogs := cowin.CentersByAge{Age18: *cd18, Age45: *cd45}
	logJSON, _ := json.MarshalIndent(centerLogs, "", " ")

	if len(*cd18) > 0 || len(*cd45) > 0 {
		msg := renderTemplate(centerLogs, vaccineDate, telegramMsgTemplate)
		notifier.SendTelegramNotification(msg, true, channelID)

		msg = fmt.Sprintf("%v - Found Available centers for %v", time.Now().Format("01-02-2006 15:04:05"), vaccineDate)
		notifier.SendTelegramNotification(msg, false, debugChannelID)

		log.Printf("Centers Available")
		_ = ioutil.WriteFile(fmt.Sprintf("%v-centers-%v.json", vaccineDate, time.Now().Format("01-02-2006-15-04-05")), logJSON, 0644)
	} else {
		msg := fmt.Sprintf("%v - No Available centers for %v", time.Now().Format("01-02-2006 15:04:05"), vaccineDate)
		notifier.SendTelegramNotification(msg, false, debugChannelID)
		log.Print("Centers Not Availabe")
	}
}

func renderTemplate(centers cowin.CentersByAge, vaccineDate string, templateString string) string {
	vaccineData := struct {
		Date           string
		TotalCenters45 int
		TotalCenters18 int
		CenterDetails  cowin.CentersByAge
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
	_ = parsedTemplate.Execute(&templateBytes, vaccineData)
	return templateBytes.String()
}

func sendEmailNotifications(centers cowin.CentersByAge, vaccineDate string, emails []string) {
	emailMsg := renderTemplate(centers, vaccineDate, emailTemplate)
	for _, mail := range emails {
		notifier.SendMail(emailMsg, "User", mail)
	}
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

func GetNextLocaleDay(locale string) string {
	loc, _ := time.LoadLocation(locale)
	now := time.Now().In(loc)
	now = now.AddDate(0, 0, 1)
	tomorrowDateIST := now.Format("02-01-2006")
	return tomorrowDateIST
}
