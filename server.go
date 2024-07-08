package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Println("No .env.local file found")
	}
}

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// NO AUTH, RATE LIMITED
	e.POST("/player_locations", func(c echo.Context) error {
		var reqBody struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		}
		if err := c.Bind(&reqBody); err != nil {
			return err
		}
		latitude := reqBody.Latitude
		longitude := reqBody.Longitude

		latitudeFloat, err := strconv.ParseFloat(latitude, 64)
		if err != nil {
			return err
		}
		longitudeFloat, err := strconv.ParseFloat(longitude, 64)
		if err != nil {
			return err
		}

		err = insertPlayerLocation(latitudeFloat, longitudeFloat)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, "Player location inserted successfully")
	})

	// AUTH
	e.GET("player_locations", func(c echo.Context) error {
		token := c.QueryParam("token")
		playerLocations, err := getAllPlayerLocations(token)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, playerLocations)
	})

	// NO AUTH
	e.GET("/current_circle", func(c echo.Context) error {
		circle, err := getCurrentCircle()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, circle)
	})

	// AUTH
	e.PUT("/set_start_conditions", func(c echo.Context) error {
		auth_token := c.Request().Header.Get("Authorization")
		startTimeStr := c.QueryParam("startTime")
		intervalStr := c.QueryParam("interval")

		startTime, err := time.Parse(standardTimeFormat, startTimeStr)
		if err != nil {
			return err
		}
		interval, err := time.ParseDuration(intervalStr)
		if err != nil {
			return err
		}
		err = populateCircles(startTime, interval, auth_token)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, "Start conditions set")
	})

	// e.GET("/test", func(c echo.Context) error {
	// 	ip := c.RealIP()
	// 	return c.String(http.StatusOK, "Client IP: "+ip)
	// })

	e.Logger.Fatal(e.Start(":1323"))
}

func getAuthToken(c echo.Context) (string, error) {
	auth_token := c.Request().Header.Get("Authorization")
	if auth_token == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	return auth_token, nil
}
