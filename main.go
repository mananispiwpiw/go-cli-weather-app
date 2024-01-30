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

	city := os.Args[1]

	// Construct the API URL using the user's city input
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=1&aqi=no&alerts=no", os.Getenv("WEATHER_API_URL_1"), city)

	res, err := http.Get(apiURL)
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

	headerMessage := fmt.Sprintf("\n%s,%s:%0.f°C UV %0.f,Condition: %s\n",
		location.Name,
		location.Country,
		current.Temp_c,
		current.Uv,
		current.Condition.Text,
	)

	coloredHeader := color.YellowString(headerMessage)

	fmt.Print(coloredHeader)

	//Header style
	hourDegreeWidth := 15
	chanceRainWidth := 15
	conditionWidth := 30
	//fmt.Println("Hour - Degree|", "\tChance Rain", "\tCondition \t\t|")

	header := fmt.Sprintf("|%-*s|%-*s|%-*s|\n",
		hourDegreeWidth, "Hour - Degree",
		chanceRainWidth, "Chance Rain",
		conditionWidth, "Condition")

	//Header Upper Border
	for i := 0; i < len(header)-1; i++ {
		fmt.Print("-")
	}
	//Craete Space
	fmt.Println()

	//Header itself
	fmt.Print(header)

	//Header Bottom Border
	for i := 0; i < len(header)-1; i++ {
		fmt.Print("-")
	}

	//Create Space
	fmt.Println("")

	for _, hour := range hours {
		date := time.Unix(int64(hour.Time_epoch), 0)

		if date.Before(time.Now()) {
			continue
		}

		//Store hour and degree as a string value
		hourDegree := fmt.Sprintf("%s - %0.f°C", date.Format("15:04"), hour.Temp_c)
		//Store chance rain value
		chanceRain := fmt.Sprintf("%0.0f%%", hour.Chance_of_rain)

		//%0.f%%
		message := fmt.Sprintf("|%-*s|%-*s|%-*s|\n",
			hourDegreeWidth, hourDegree,
			chanceRainWidth, chanceRain,
			conditionWidth, hour.Condition.Text,
		)

		if hour.Chance_of_rain < 45 {
			fmt.Print(message)
		} else if hour.Chance_of_rain <= 60 {
			color.Cyan(message)
		} else {
			color.Blue(message)
		}
	}

	//Table Bottom Border
	for i := 0; i < len(header)-1; i++ {
		fmt.Print("-")
	}
	//Create Space
	fmt.Println()

}
