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

// TODO: rename
type responseStruct struct {
	Expand struct {
		CurrentCircle Circle `json:"current_circle"`
	} `json:"expand"`
}

type PlayerLocation struct {
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	ID             string `json:"id"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Updated        string `json:"updated"`
}

type PlayerLocationsResponse struct {
	Page       int              `json:"page"`
	PerPage    int              `json:"perPage"`
	TotalItems int              `json:"totalItems"`
	TotalPages int              `json:"totalPages"`
	Items      []PlayerLocation `json:"items"`
}
