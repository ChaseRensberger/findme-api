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
	fmt.Println("-----------------")
	fmt.Printf("Start: %s\n", response.Items[0].Start)
	fmt.Printf("End: %s\n", response.Items[0].End)
	fmt.Printf("Width: %d\n", response.Items[0].Meters)
	fmt.Println("-----------------")
	updateCircle(response.Items[0], auth_token)

	// now update every circle based on the previous circle
	for circleIdx := 1; circleIdx < len(response.Items); circleIdx++ {
		currentCircle := response.Items[circleIdx]
		currentCircle.Start = response.Items[circleIdx-1].End
		end, err := time.Parse(standardTimeFormat, response.Items[circleIdx-1].End)
		if err != nil {
			return err
		}
		currentCircle.End = end.Add(interval).Format(standardTimeFormat)
		updateCircle(currentCircle, auth_token)
	}

	return nil
}

func getCurrentCircle() (Circle, error) {
	apiUrl := os.Getenv("POCKETBASE_URL")
	req, err := http.NewRequest(http.MethodGet, apiUrl+"/api/collections/circles/records", nil)
	if err != nil {
		return Circle{}, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return Circle{}, err
	}

	currentTime := time.Now()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Circle{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response CirclesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Circle{}, err
	}

	// Find the circle where start <= currentTime < end
	for _, circle := range response.Items {

		start, err := time.Parse(standardTimeFormat, circle.Start)
		if err != nil {
			fmt.Println(err)
			return Circle{}, err
		}
		fmt.Println(start)

		end, err := time.Parse(standardTimeFormat, circle.End)
		if err != nil {
			return Circle{}, err
		}

		if currentTime.After(start) && currentTime.Before(end) {
			return circle, nil
		}
	}

	return Circle{}, fmt.Errorf("no circle found for current time")
}

func authorizeRequest(req *http.Request, auth_token string) {
	req.Header.Set("Authorization", "Bearer "+auth_token)
}
