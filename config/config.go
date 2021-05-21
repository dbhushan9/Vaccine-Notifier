package config

type NotificationOptions struct {
	DebugTelegramChannels []string `json:"debug_channel_ids"`
	TelegramChannels      []string `json:"channel_ids"`
	Emails                []string `json:"emails"`
}

type CenterOptions struct {
	Date              string `json:"vaccine_data"`
	DistrictId        int    `json:"district_id"`
	BlockName         string `json:"block_name"`
	Vaccine           string `json:"vaccine"`
	FeeType           string `json:"fee_type"`
	AgeSlots          []int  `json:"age_slots"`
	ExcludedCenterIds []int  `json:"excluded_center_ids"`
}

type VaccineQuery struct {
	CenterOptions       CenterOptions       `json:"center_options"`
	NotificationOptions NotificationOptions `json:"notification_options"`
}

type VaccineQueries = []VaccineQuery
