package cowin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func QueryCowinAPI(vaccineDate string, districtId int) (*CowinAPIResponse, error) {
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
		log.Print("Error reading response body")
		log.Print(err.Error())
	}

	var response CowinAPIResponse
	err = json.Unmarshal(bodyBytes, &response)
	return &response, err
}
