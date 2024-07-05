package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func insertPlayerLocation(latitude float64, longitude float64) error {
	payload := map[string]float64{
		"latitude":  latitude,
		"longitude": longitude,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	apiUrl := os.Getenv("POCKETBASE_URL")
  // why able to pass in something not of io.Reader if we don't specify type
	req, err := http.NewRequest(http.MethodPost, apiUrl+"/api/collections/player_locations/records", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// TODO
func getAllPlayerLocations(auth_token) {

	apiUrl := os.Getenv("POCKETBASE_URL")
	req, err := http.NewRequest(http.MethodGet, apiUrl+"/api/collections/player_locations/records", nil)
	if err != nil {
		return err
	}
  authorizeRequest(req, auth_token)
	client := http.DefaultClient
	resp, err := client.Do(req)
  //...

}

func getCurrentCircleByGameId(id string, auth_token string) (Circle, error) {
	apiUrl := os.Getenv("POCKETBASE_URL")
	req, err := http.NewRequest(http.MethodGet, apiUrl+"/api/collections/games/records/"+id+"?expand=current_circle", nil)
	if err != nil {
		return Circle{}, err
	}
  authorizeRequest(req, auth_token)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return Circle{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Circle{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response responseStruct
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Circle{}, err
	}

	circle := response.Expand.CurrentCircle
	return circle, nil

}

// TODO: Verify works
func authorizeRequest(req, auth_token) {
	req.Header.Set("Authorization", "Bearer " + auth_token)
}
