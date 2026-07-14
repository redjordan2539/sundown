package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Data struct {
	SunsetTime           time.Time `json:"sunset"`
	TenMinuteWarning     time.Time `json:"ten_minute_warning"`
	LastUpdated          time.Time `json:"last_updated"`
	MorningNotifySent    bool      `json:"morning_notify_sent"`
	TenMinuteWarningSent bool      `json:"ten_minute_warning_sent"`
}

func NewData() *Data {
	return &Data{
		MorningNotifySent:    false,
		TenMinuteWarningSent: false,
	}
}

func loadData(path string) (*Data, error) {
	state, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewData(), nil
		} else {
			return nil, err
		}
	}
	if len(state) == 0 {
		return NewData(), nil
	}

	var d Data
	err = json.Unmarshal(state, &d)
	if err != nil {
		return nil, err
	}
	return &d, err
}

func (d *Data) Save(path string) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	tmpPath := path + ".path"
	defer os.Remove(tmpPath)

	if err := os.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	return nil
}

func (d *Data) HandleSunsetNotifications(path string) error {
	if !d.MorningNotifySent {
		t := d.SunsetTime.Format("15:04")
		err := SendNotification(fmt.Sprintf("Sunset is at %s", t))
		if err != nil {
			return fmt.Errorf("Unable to send notification | %s", err)
		}
		d.MorningNotifySent = true
		d.Save(path)
	}

	if time.Now().After(d.TenMinuteWarning) && !d.TenMinuteWarningSent {
		err := SendNotification("Sunset is Approaching!")
		if err != nil {
			return fmt.Errorf("Unable to send notification | %s", err)
		}
		d.TenMinuteWarningSent = true
		d.Save(path)
	}
	return nil
}
