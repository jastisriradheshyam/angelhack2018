package main

import (
"time"

)

type bottle struct {

	BottleId string `json:"bottleId"`
	BloodGroup string `json:"bloodGroup"`
	CurrentOwner string `json:"currentOwner"`
	DateOfPacking time.Time `json:"dateOfPacking"`
	Trail string `json:"trail"`
	Status string `json:"status"`
	SentTo string `json:sentTo,omitEmpty`


}

type user struct {
	UserId string `json:"userId"`
	CurrentStock map[string]stockRequirement `json:"currentStock"`
	CurrentRequirement map[string]int `json:"currentRequirement,omitempty"`
	EmailId string `json:"emailId"`
	ContactPerson string `json:"currentPerson"`
	Region string `json:"region"`
	Type string `json:"type"`
	Asking []string `json:"Asking"`
	Giving  []string `json:"Giving"`
	
}
	
	type stockRequirement struct {
		BloodGroup  string `json:"bloodGroup"`
		BottleMap []string `json:"bottleMap"`
		Count int `json:"count"`
		MustCount int `json:"mustCount",omitempty`
	}

