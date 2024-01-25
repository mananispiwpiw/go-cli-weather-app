package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	res, err := http.Get("http://api.weatherapi.com/v1/current.json?key=c01fc842f0664034bda42040242501&q=Jakarta&aqi=no")
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

	fmt.Println(string(body))
}
