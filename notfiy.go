package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func SendNotification(m string) error {
	req, err := http.NewRequest(http.MethodPost, "https://ntfy.delpilar.net/test", strings.NewReader(m))
	if err != nil {
		return err
	}

	ntfyUser := os.Getenv("NTFY_USER")
	ntfyPass := os.Getenv("NTFY_PASS")

	if ntfyUser == "" || ntfyPass == "" {
		return fmt.Errorf("ntfy username or password not set in environment variables. Please set them and try again")
	}

	req.SetBasicAuth(ntfyUser, ntfyPass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
