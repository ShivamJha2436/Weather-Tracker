package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"openWeatherMapApiKey"`
}

type weatherData struct{
	Name string `json:"name"`
	Main struct{
		Kelvin float64 `json:"temp"`
	}`json:"main"`
}

func loadApiCongig(filename string) (apiConfigData, error){
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var config apiConfigData

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return apiConfigData{}, err
	}
	return config, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func query(city string, w http.ResponseWriter) (weatherData, error){
	apiConfig, err := loadApiCongig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)

	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData 
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}


func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
				city := strings.SplitN(r.URL.Path, "/", 3)[2]
				data, err := query(city, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(data)
		})	

	http.ListenAndServe(":8080", nil)
}