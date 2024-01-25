package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	fmt.Println(weather)
}
