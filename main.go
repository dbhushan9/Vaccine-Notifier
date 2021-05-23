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
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
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
	log.SetReportCaller(true)

	log.SetOutput(&lumberjack.Logger{
		Filename:   "/app/shared/out.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     2,    //days
		Compress:   true, // disabled by default
	})
}

func main() {
	log.Info("starting vaccine alert worker")
	date := GetNextLocaleDay(indianLocale)
	ageSlots2 := []int{18, 45}
	districtId := 363
	blockName := "Haveli"
	vaccine := "any"
	feeType := "any"
	//ARMY ONLY CENTERS
	excludedCenter := []int{619964, 629727}

	channelID := os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT")
	debugChannelID := os.Getenv("TELEGRAM_CHANNEL_ID_VACCINE_ALERT_DEBUG")
	var emails = []string{"dbhushan912@gmail.com"}

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
		defer ExecutionTimeTaken(time.Now(), "cronjob")
		centersByAge := checkForVaccineCenters(query.CenterOptions)
		sendNotification(query.NotificationOptions, *centersByAge, query.CenterOptions.Date)
	})
	<-gocron.Start()
}

//TODO: send error if failed inbetween and send debug mail accordingly
func checkForVaccineCenters(centerOptions config.CenterOptions) *cowin.CentersByAge {
	defer ExecutionTimeTaken(time.Now(), "checkForVaccineCenters")
	log.WithFields(log.Fields{"center_criteria": centerOptions}).Info("checking for vaccine centers")

	response, err := cowin.QueryCowinAPI(centerOptions.Date, centerOptions.DistrictId)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error unmarshalling request body")
	}

	log.WithFields(log.Fields{"center_count": len(response.Centers)}).Info("found centers")

	//TODO: call process only once, return response by groupBy
	_, cd18 := cowin.ProcessCentersPresent(response.Centers, centerOptions.AgeSlots[0], centerOptions.Vaccine, centerOptions.FeeType, centerOptions.Date, centerOptions.BlockName, centerOptions.ExcludedCenterIds)
	_, cd45 := cowin.ProcessCentersPresent(response.Centers, centerOptions.AgeSlots[1], centerOptions.Vaccine, centerOptions.FeeType, centerOptions.Date, centerOptions.BlockName, centerOptions.ExcludedCenterIds)

	centersByAge := cowin.CentersByAge{Age18: *cd18, Age45: *cd45}

	return &centersByAge

}

func sendNotification(options config.NotificationOptions, centersByAge cowin.CentersByAge, date string) {
	defer ExecutionTimeTaken(time.Now(), "sendNotification")
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

		msg = fmt.Sprintf("%v - Found available centers for %v", getCurrentTime(), date)
		for _, debugChannelID := range options.DebugTelegramChannels {
			notifier.SendTelegramNotification(msg, false, debugChannelID)
		}
		log.Info("found vaccine centers")

	} else {
		msg := fmt.Sprintf("%v - No available centers for %v", getCurrentTime(), date)
		for _, debugChannelID := range options.DebugTelegramChannels {
			notifier.SendTelegramNotification(msg, false, debugChannelID)
		}
		log.Info("no vaccine centers found")
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

func GetNextLocaleDay(locale string) string {
	loc, _ := time.LoadLocation(locale)
	now := time.Now().In(loc)
	now = now.AddDate(0, 0, 1)
	tomorrowDateIST := now.Format("02-01-2006")
	return tomorrowDateIST
}

func getCurrentTime() string {
	return time.Now().Format("01-02-2006 15:04:05")
}

func ExecutionTimeTaken(t time.Time, n string) {
	e := time.Since(t)
	log.WithFields(log.Fields{
		"execution_time": e,
		"function_name":  n,
	}).Info("measured execution time")
}
