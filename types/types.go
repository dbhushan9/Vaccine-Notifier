package types

import (
	arrayUtils "dbhushan9/vaccine-alerts/util"
	"log"
)

//ARMY ONLY CENTERS
var excludedCenter = []int{619964, 629727}

type ReponseType struct {
	Centers SessionCalendarEntriesSchema `json:"centers"`
}

// SessionCalendarEntriesSchema defines model for SessionCalendarEntriesSchema.
type SessionCalendarEntriesSchema []SessionCalendarEntrySchema

type SessionCalendarEntrySchema struct {
	Address *string `json:"address,omitempty"`

	// Address line in preferred language as specified in Accept-Language header parameter.
	AddressL  *string `json:"address_l,omitempty"`
	BlockName string  `json:"block_name"`

	// Block name in preferred language as specified in Accept-Language header parameter.
	BlockNameL   *string `json:"block_name_l,omitempty"`
	CenterId     float32 `json:"center_id"`
	DistrictName string  `json:"district_name"`

	// District name in preferred language as specified in Accept-Language header parameter.
	DistrictNameL *string `json:"district_name_l,omitempty"`

	// Fee charged for vaccination
	FeeType SessionCalendarEntrySchemaFeeType `json:"fee_type"`
	From    string                            `json:"from"`
	Lat     *float32                          `json:"lat,omitempty"`
	Long    *float32                          `json:"long,omitempty"`
	Name    string                            `json:"name"`

	// Name in preferred language as specified in Accept-Language header parameter.
	NameL     *string   `json:"name_l,omitempty"`
	Pincode   int       `json:"pincode"`
	Sessions  []Session `json:"sessions"`
	StateName string    `json:"state_name"`

	// State name in preferred language as specified in Accept-Language header parameter.
	StateNameL  *string               `json:"state_name_l,omitempty"`
	To          string                `json:"to"`
	VaccineFees *VaccineFeeListSchema `json:"vaccine_fees,omitempty"`
}

type Session struct {
	AvailableCapacity float32 `json:"available_capacity"`
	Date              string  `json:"date"`
	MinAgeLimit       int     `json:"min_age_limit"`
	SessionId         string  `json:"session_id"`

	// Array of slot names
	Slots   []string `json:"slots"`
	Vaccine string   `json:"vaccine"`
}

// Fee charged for vaccination
type SessionCalendarEntrySchemaFeeType string

// VaccineFeeListSchema defines model for VaccineFeeListSchema.
type VaccineFeeListSchema []VaccineFeeSchema

type VaccineFeeSchema struct {
	Fee     string `json:"fee"`
	Vaccine string `json:"vaccine"`
}

type LogDetails struct {
	Age18 []CenterDetails `json:"age18"`
	Age45 []CenterDetails `json:"age45"`
}

type CenterDetails struct {
	Name              string  `json:"centerName"`
	Date              string  `json:"date"`
	AvailableCapacity float32 `json:"availableCapacity"`
	FeeType           string  `json:"fees"`
	Vaccine           string  `json:"vaccine"`
}

func ProcessCentersPresent(centers SessionCalendarEntriesSchema, ageSlot int, vaccine string, feeType string, date string, blockName string) (*SessionCalendarEntriesSchema, *[]CenterDetails) {
	validCenters := make(SessionCalendarEntriesSchema, 0)
	centerData := make([]CenterDetails, 0)

	for _, center := range centers {
		if center.isValidCenter(ageSlot, vaccine, feeType, date, blockName) {
			validCenters = append(validCenters, center)
			session := center.Sessions[0]
			centerData = append(centerData, CenterDetails{Name: center.Name, Date: session.Date, AvailableCapacity: session.AvailableCapacity, FeeType: string(center.FeeType), Vaccine: session.Vaccine})
		}
	}
	return &validCenters, &centerData
}

func (center *SessionCalendarEntrySchema) isValidCenter(ageSlot int, vaccine string, feeType string, date string, blockName string) bool {
	if (blockName != "any" && center.BlockName != blockName) || arrayUtils.Includes(excludedCenter, int(center.CenterId)) {
		log.Printf("Skipping center %v %v %v", center.CenterId, center.Name, center.BlockName)
		return false
	}
	validSessions := []Session{}
	for _, session := range center.Sessions {
		if session.isValidSession(ageSlot, vaccine, date) {
			validSessions = append(validSessions, session)
		}
	}
	return len(validSessions) > 0
}

func (s *Session) isValidSession(ageSlot int, vaccine string, date string) bool {
	return int(s.AvailableCapacity) > 0 && s.Date == date && s.MinAgeLimit == ageSlot && (vaccine == "any" || s.Vaccine == vaccine)
}
