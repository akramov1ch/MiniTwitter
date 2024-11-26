package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const apiURL = "https://api.octopush.com/v1/public/multi-channel/send"

func SendSMS(phone, message string) error {
	data := map[string]interface{}{
		"channel": "sms",
		"text":    message,
		"recipients": []map[string]string{
			{"phone_number": phone},
		},
		"sender":             "NewMsg",
		"auto_optimize_text": true,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-login", "shaxbozakramovic@gmail.com")
	req.Header.Set("api-key", "rNXM1IPAL26UKkB5pyT0sWmtSzdJYcZF")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send SMS: %s", resp.Status)
	}

	return nil
}
