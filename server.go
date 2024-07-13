package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Println("No .env.local file found")
	}
}

func main() {

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                        // Allows all origins
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}, // Specify allowed methods
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// NO AUTH, RATE LIMITED
	e.POST("/player_location", func(c echo.Context) error {
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
		circle := getCurrentCircle()
		if circle.ID == "" {
			return c.JSON(http.StatusNotFound, struct{}{})
		}
		return c.JSON(http.StatusOK, circle)
	})

	// AUTH
	e.PUT("/start_conditions", func(c echo.Context) error {
		auth_token, err := getAuthToken(c)
		if err != nil {
			return err
		}
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

	e.GET("/game_state", func(c echo.Context) error {
		id := c.QueryParam("id")
		gameState, err := getGameState(id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, gameState)
	})

	// e.GET("/test", func(c echo.Context) error {
	// 	ip := c.RealIP()
	// 	return c.String(http.StatusOK, "Client IP: "+ip)
	// })

	e.Logger.Fatal(e.Start(":1324"))
}

func getAuthToken(c echo.Context) (string, error) {
	auth_token := c.Request().Header.Get("Authorization")
	if auth_token == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	auth_token = strings.TrimPrefix(auth_token, "Bearer ")
	return auth_token, nil
}
