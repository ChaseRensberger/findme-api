package main

type Circle struct {
	Latitude       string `json:"Latitude"`
	Longitude      string `json:"Longitude"`
	Meters         int    `json:"Meters"`
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Field          string `json:"field"`
	ID             string `json:"id"`
	Updated        string `json:"updated"`
}

type responseStruct struct {
	Expand struct {
		CurrentCircle Circle `json:"current_circle"`
	} `json:"expand"`
}
