package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		TempFeels float64 `json:"feelslike_c"`
		Wind      float64 `json:"wind_kph"`
		Humidity  float64 `json:"humidity"`

		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				Timepoch  int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				RainChance float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Baku"
	api := "Your API Key"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + api + "" + q + "&days=3&aqi=no&alerts=yes")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API is not available or you did not provide one. Please add it to the source code.")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"| %s, %s / %.0f celsium, but feels like %.0f / %s / Wind Speed: %.0fkmph / Humidity: %.0f |\n",
		location.Name,
		location.Country,
		current.TempC,
		current.TempFeels,
		current.Condition.Text,
		current.Wind,
		current.Humidity,
	)

	for _, hour := range hours {
		date := time.Unix(hour.Timepoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s - %.0f Celsium, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.RainChance,
			hour.Condition.Text,
		)

		if hour.RainChance < 40 {
			fmt.Print(message)
		}

		if hour.RainChance < 65 && hour.RainChance > 40 {
			color.Yellow(message)
		} else {
			color.Red(message)
		}
	}
}
