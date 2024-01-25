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

// Struct to store data from API
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct { //Current status
		Temp_c    float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		}
		Uv float64 `json:"uv"`
	} `json:"current"`
	Forecast struct { //Forecasting or prediction
		Forecasatday []struct {
			Hour []struct {
				Time_epoch     int     `json:"time_epoch"`
				Temp_c         float64 `json:"temp_c"`
				Chance_of_rain float64 `json:"chance_of_rain"`
				Condition      struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"hour"`
		} `json:"forecastday"`
	}
}

func main() {

	res, err := http.Get(os.Getenv("WEATHER_API_URL_1"))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available!")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	//Store the data from WeatherAPI to struct
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, hours, current := weather.Location, weather.Forecast.Forecasatday[0].Hour, weather.Current

	headerMessage := fmt.Sprintf("%s,%s:%0.f°C,%s\n",
		location.Name,
		location.Country,
		current.Temp_c,
		current.Condition.Text,
	)

	coloredHeader := color.YellowString(headerMessage)

	fmt.Print(coloredHeader)
	fmt.Println("Hour - Degree", "\tChance Rain", "\tCondition")

	for _, hour := range hours {
		date := time.Unix(int64(hour.Time_epoch), 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %0.f°C,\t%0.f%%,\t\t%s\n",
			date.Format("15:04"),
			hour.Temp_c,
			hour.Chance_of_rain,
			hour.Condition.Text,
		)

		if hour.Chance_of_rain < 45 {
			fmt.Print(message)
		} else if hour.Chance_of_rain <= 60 {
			color.Cyan(message)
		} else {
			color.Blue(message)
		}
	}
}
