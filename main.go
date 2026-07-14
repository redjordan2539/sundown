package main

import (
	"log"
	"log/slog"
	"os"
	"time"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

func main() {
	sunsetJSONPath := os.Getenv("SUNSET_JSON_PATH")
	sunset, err := loadData(sunsetJSONPath)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if !isToday(sunset.LastUpdated) {
		slog.Info("Sunset Data out of date. Refreshing.", "last_update", sunset.LastUpdated)

		lat := os.Getenv("LAT")
		lng := os.Getenv("LNG")
		apiResponse, err := getSunsetData(lat, lng)
		if err != nil {
			slog.Error("Unable to get Sunset data from API", "err", err)
			return
		}

		sunset = &Data{
			SunsetTime:           apiResponse.SunsetTime,
			TenMinuteWarning:     apiResponse.SunsetTime.Add(-10 * time.Minute),
			LastUpdated:          time.Now(),
			MorningNotifySent:    false,
			TenMinuteWarningSent: false,
		}
		sunset.Save(sunsetJSONPath)
		slog.Info("Sunset Data Saved", "path", sunsetJSONPath)
	}

	sunset.HandleSunsetNotifications(sunsetJSONPath)
}
