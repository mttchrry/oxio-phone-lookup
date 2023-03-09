package model

type PhoneNumber struct{
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
	AreaCode    string `json:"areaCode"`
	Local  	    string `json:"localPhoneNumber"`
}


