package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"dbhushan9/vaccine-alerts/config"
	"dbhushan9/vaccine-alerts/cowin"
	"dbhushan9/vaccine-alerts/notifier"

	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
)

const (
	indianLocale = "Asia/Kolkata"
)

//go:embed "templates/telegram-notification.md"
var telegramMsgTemplate string

//go:embed "templates/email-template.html"
var emailTemplate string

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "@timestamp",
			log.FieldKeyMsg:  "message",
		},
	})
	// err := os.Mkdir("shared", 0666)
	// log.WithFields(log.Fields{"error": err}).Error("failed to create shared directory")
	file, err := os.OpenFile("/app/shared/out.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to create log file")

	}
	log.SetOutput(file)

	// defer file.Close()
	log.SetOutput(file)
	log.Info("completed logging init")
}

func main() {
	// setLoggingConfig()
	log.Info("Starting Vaccine alert worker")
	date := GetNextLocaleDay(indianLocale)
	// ageSlots := [2]int{18, 45}
	ageSlots2 := []int{18, 45}
	districtId := 363
	blockName := "Haveli"
	vaccine := "any"
	feeType := "any"
	//ARMY ONLY CENTERS
	excludedCenter := []int{619964, 629727}

	channelID := os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT")
	debugChannelID := os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT_DEBUG")
	var emails = []string{"dbhushan912@gmail.com", "avi6387@gmail.com", "deshmukh.kalyani81@gmail.com", "rupadeshmukh26@gmail.com"}

	query := &config.VaccineQuery{
		CenterOptions: config.CenterOptions{
			Date:              date,
			DistrictId:        districtId,
			BlockName:         blockName,
			Vaccine:           vaccine,
			FeeType:           feeType,
			AgeSlots:          ageSlots2,
			ExcludedCenterIds: excludedCenter,
		},
		NotificationOptions: config.NotificationOptions{
			DebugTelegramChannels: []string{debugChannelID},
			TelegramChannels:      []string{channelID},
			Emails:                emails,
		},
	}

	gocron.Every(5).Minute().Do(func() {
		centersByAge := checkForVaccineCenters(query.CenterOptions)
		sendNotification(query.NotificationOptions, *centersByAge, query.CenterOptions.Date)
	})
	<-gocron.Start()
	// checkForVaccineCenters(vaccineDate, districtId)
}

//TODO: send error if failed inbetween and send debug mail accordingly
func checkForVaccineCenters(centerOptions config.CenterOptions) *cowin.CentersByAge {
	log.WithFields(log.Fields{"centerOptions": centerOptions}).Info("checking for vaccine centers")

	response, err := cowin.QueryCowinAPI(centerOptions.Date, centerOptions.DistrictId)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error unmarshalling request body")
	}

	log.Info("Total centers are", len(response.Centers))

	//TODO: call process only once, return response by groupBy
	_, cd18 := cowin.ProcessCentersPresent(response.Centers, centerOptions.AgeSlots[0], centerOptions.Vaccine, centerOptions.FeeType, centerOptions.Date, centerOptions.BlockName, centerOptions.ExcludedCenterIds)
	_, cd45 := cowin.ProcessCentersPresent(response.Centers, centerOptions.AgeSlots[1], centerOptions.Vaccine, centerOptions.FeeType, centerOptions.Date, centerOptions.BlockName, centerOptions.ExcludedCenterIds)

	centersByAge := cowin.CentersByAge{Age18: *cd18, Age45: *cd45}

	return &centersByAge

}

func sendNotification(options config.NotificationOptions, centersByAge cowin.CentersByAge, date string) {
	//TODO:
	//Send Telegram Success Message
	//Send Telegram Debug success Message
	//Send Sucess Email

	//Send Telegram Debug failure Message

	if len(centersByAge.Age18) > 0 || len(centersByAge.Age45) > 0 {
		//TODO: send a single message with grouped values
		msg := renderTemplate(centersByAge, date, telegramMsgTemplate)
		for _, channelID := range options.TelegramChannels {
			notifier.SendTelegramNotification(msg, true, channelID)
		}

		msg = fmt.Sprintf("%v - Found Available centers for %v", time.Now().Format("01-02-2006 15:04:05"), date)
		for _, debugChannelID := range options.DebugTelegramChannels {
			notifier.SendTelegramNotification(msg, false, debugChannelID)
		}
		log.Info("Centers Available")

	} else {
		msg := fmt.Sprintf("%v - No Available centers for %v", time.Now().Format("01-02-2006 15:04:05"), date)
		for _, debugChannelID := range options.DebugTelegramChannels {
			notifier.SendTelegramNotification(msg, false, debugChannelID)
		}
		log.Info("Centers Not Availabe")
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

}

func GetNextLocaleDay(locale string) string {
	loc, _ := time.LoadLocation(locale)
	now := time.Now().In(loc)
	now = now.AddDate(0, 0, 1)
	tomorrowDateIST := now.Format("02-01-2006")
	return tomorrowDateIST
}
