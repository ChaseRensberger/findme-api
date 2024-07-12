package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"
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

func getAllPlayerLocations(auth_token string) ([]PlayerLocation, error) {
	apiUrl := os.Getenv("POCKETBASE_URL")
	allLocations := []PlayerLocation{}

	page := 1
	for {
		reqUrl := fmt.Sprintf("%s/api/collections/player_locations/records?page=%d", apiUrl, page)
		req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		if err != nil {
			return nil, err
		}
		authorizeRequest(req, auth_token)

		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var response PlayerLocationsResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}

		allLocations = append(allLocations, response.Items...)

		if page >= response.TotalPages {
			break
		}
		page++
	}

	return allLocations, nil
}

func updateCircle(circle Circle, auth_token string) error {
	apiUrl := os.Getenv("POCKETBASE_URL")
	payload, err := json.Marshal(circle)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, apiUrl+"/api/collections/circles/records/"+circle.ID, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	authorizeRequest(req, auth_token)
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

// Takes a circle and a start time and populates every following circle based on decreasing width with corresponding start and end times based on a provided interval
func populateCircles(startTime time.Time, interval time.Duration, auth_token string) error {
	// first get all circles
	apiUrl := os.Getenv("POCKETBASE_URL")
	req, err := http.NewRequest(http.MethodGet, apiUrl+"/api/collections/circles/records", nil)
	if err != nil {
		return err
	}
	authorizeRequest(req, auth_token)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response CirclesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	// now sort the circles by width
	sort.Slice(response.Items, func(i, j int) bool {
		return response.Items[i].Meters > response.Items[j].Meters
	})

	// update first circle
	response.Items[0].Start = startTime.Format(standardTimeFormat)
	response.Items[0].End = startTime.Add(interval).Format(standardTimeFormat)
	updateCircle(response.Items[0], auth_token)

	// now update every circle based on the previous circle
	for circleIdx := 1; circleIdx < len(response.Items); circleIdx++ {
		response.Items[circleIdx].Start = response.Items[circleIdx-1].End
		end, err := time.Parse(standardTimeFormat, response.Items[circleIdx-1].End)
		if err != nil {
			return err
		}
		response.Items[circleIdx].End = end.Add(interval).Format(standardTimeFormat)
		updateCircle(response.Items[circleIdx], auth_token)
	}

	setGameTime("qji7yh0msqphiq9", response.Items[0].Start, response.Items[len(response.Items)-1].End, auth_token)

	return nil
}

func setGameTime(id string, start string, end string, auth_token string) error {
	apiUrl := os.Getenv("POCKETBASE_URL")
	payload := map[string]string{
		"start": start,
		"end":   end,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, apiUrl+"/api/collections/games/records/"+id, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	authorizeRequest(req, auth_token)

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

func getCurrentCircle() Circle {
	apiUrl := os.Getenv("POCKETBASE_URL")
	auth_token := os.Getenv("AUTH_TOKEN")
	req, err := http.NewRequest(http.MethodGet, apiUrl+"/api/collections/circles/records", nil)
	if err != nil {
		return Circle{}
	}
	authorizeRequest(req, auth_token)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return Circle{}
	}

	currentTime := time.Now().UTC()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Circle{}
	}

	var response CirclesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Circle{}
	}

	// Find the circle where start <= currentTime < end
	for _, circle := range response.Items {

		start, err := time.Parse(standardTimeFormat, circle.Start)
		if err != nil {
			fmt.Println(err)
			return Circle{}
		}

		end, err := time.Parse(standardTimeFormat, circle.End)
		if err != nil {
			return Circle{}
		}

		if currentTime.After(start) && currentTime.Before(end) {
			return circle
		}
	}

	return Circle{}
}

func getGameState(id string) (GameState, error) {
	apiUrl := os.Getenv("POCKETBASE_URL")
	req, err := http.NewRequest(http.MethodGet, apiUrl+"/api/collections/games/records/"+id, nil)
	if err != nil {
		return "", err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var game Game
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		return "", err
	}

	start, err := time.Parse(standardTimeFormat, game.Start)
	if err != nil {
		return "", err
	}

	end, err := time.Parse(standardTimeFormat, game.End)
	if err != nil {
		return "", err
	}

	currentTime := time.Now().UTC()

	if currentTime.Before(start) {
		return WAITING, nil
	}

	if currentTime.After(start) && currentTime.Before(end) {
		return ACTIVE, nil
	}

	return FINISHED, nil
}

func authorizeRequest(req *http.Request, auth_token string) {
	req.Header.Set("Authorization", "Bearer "+auth_token)
}
