package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ResponseData struct {
	SunsetTime time.Time `json:"sunset"`
}

func getSunsetData(lat string, lng string) (*ResponseData, error) {
	u, _ := url.Parse("https://api.sunrise-sunset.org/v2")

	q := u.Query()
	q.Set("lat", lat)
	q.Set("lng", lng)
	q.Set("formatted", "0")

	u.RawQuery = q.Encode()

	URL := u.String()

	r, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", r.StatusCode)
	}

	var apiResponse struct {
		SunsetTime string `json:"sunset"`
	}

	if err := json.NewDecoder(r.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	parsedTime, err := time.Parse(time.RFC3339, apiResponse.SunsetTime)
	if err != nil {
		return nil, err
	}
	return &ResponseData{SunsetTime: parsedTime}, nil
}
