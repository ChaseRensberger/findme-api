package main

var standardTimeFormat = "2006-01-02 15:04:05.000Z"

type Game struct {
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	ID             string `json:"id"`
	Updated        string `json:"updated"`
	Start          string `json:"start"`
	End            string `json:"end"`
}

type Circle struct {
	Latitude       string  `json:"latitude"`
	Longitude      string  `json:"longitude"`
	Meters         int     `json:"meters"`
	CollectionID   string  `json:"collectionId"`
	CollectionName string  `json:"collectionName"`
	Created        string  `json:"created"`
	Field          string  `json:"field"`
	ID             string  `json:"id"`
	Updated        string  `json:"updated"`
	Start          string  `json:"start"`
	End            string  `json:"end"`
	Zoom           float64 `json:"zoom"`
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

type CirclesResponse struct {
	Page       int      `json:"page"`
	PerPage    int      `json:"perPage"`
	TotalPages int      `json:"totalPages"`
	TotalItems int      `json:"totalItems"`
	Items      []Circle `json:"items"`
}

type GameState string

const (
	WAITING  GameState = "WAITING"
	ACTIVE   GameState = "ACTIVE"
	FINISHED GameState = "FINISHED"
)
