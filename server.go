package main

import (
	"fmt"
	"net/http"
	"strconv"

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

	e.POST("/insert_player", func(c echo.Context) error {
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

	e.GET("/get_current_circle", func(c echo.Context) error {
		id := c.QueryParam("id")
		token := c.QueryParam("token")

		circle, err := getCurrentCircleByGameId(id, token)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, circle)
	})

  e.GET("/test", func(c echo.Context) error {
    ip := c.RealIP()
    return c.String(http.StatusOK, "Client IP: " + ip)
  })

	e.Logger.Fatal(e.Start(":1323"))
}
