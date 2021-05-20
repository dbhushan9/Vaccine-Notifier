package cowin

type CowinAPIResponse struct {
	Centers []VaccineCenter `json:"centers"`
}

type VaccineCenters = []VaccineCenter

type VaccineCenter struct {
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
	FeeType string   `json:"fee_type"`
	From    string   `json:"from"`
	Lat     *float32 `json:"lat,omitempty"`
	Long    *float32 `json:"long,omitempty"`
	Name    string   `json:"name"`

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

// VaccineFeeListSchema defines model for VaccineFeeListSchema.
type VaccineFeeListSchema []VaccineFeeSchema

type VaccineFeeSchema struct {
	Fee     string `json:"fee"`
	Vaccine string `json:"vaccine"`
}
