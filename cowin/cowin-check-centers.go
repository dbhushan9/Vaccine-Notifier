package cowin

import log "github.com/sirupsen/logrus"

type CentersByAge struct {
	Age18 []CenterDetails `json:"age18"`
	Age45 []CenterDetails `json:"age45"`
}

type CenterDetails struct {
	Name              string  `json:"center_name"`
	Pincode           int     `json:"district_id"`
	Date              string  `json:"date"`
	AvailableCapacity float32 `json:"available_capacity"`
	FeeType           string  `json:"fees"`
	Vaccine           string  `json:"vaccine"`
}

func ProcessCentersPresent(centers VaccineCenters, ageSlot int, vaccine string, feeType string, date string, blockName string, excludedCenter []int) (*VaccineCenters, *[]CenterDetails) {
	validCenters := make(VaccineCenters, 0)
	centerData := make([]CenterDetails, 0)

	for _, center := range centers {
		if isValidCenter(center, ageSlot, vaccine, feeType, date, blockName, excludedCenter) {
			validCenters = append(validCenters, center)
			session := center.Sessions[0]
			centerData = append(centerData, CenterDetails{Name: center.Name, Pincode: center.Pincode, Date: session.Date, AvailableCapacity: session.AvailableCapacity, FeeType: string(center.FeeType), Vaccine: session.Vaccine})
		}
	}
	return &validCenters, &centerData
}

func isValidCenter(center VaccineCenter, ageSlot int, vaccine string, feeType string, date string, blockName string, excludedCenter []int) bool {
	if (blockName != "any" && center.BlockName != blockName) || includes(excludedCenter, int(center.CenterId)) {
		log.WithFields(log.Fields{"center": center}).Info("Skipping center")
		return false
	}
	validSessions := []Session{}
	for _, session := range center.Sessions {
		if isValidSession(session, ageSlot, vaccine, date) {
			validSessions = append(validSessions, session)
		}
	}
	return len(validSessions) > 0
}

func isValidSession(s Session, ageSlot int, vaccine string, date string) bool {
	return int(s.AvailableCapacity) > 0 && s.Date == date && s.MinAgeLimit == ageSlot && (vaccine == "any" || s.Vaccine == vaccine)
}

func includes(arr []int, val int) bool {
	for _, e := range arr {
		if e == val {
			return true
		}
	}
	return false
}
