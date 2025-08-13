package external

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func (w *WAApi) SendMessage(phoneNumber, message string) (err error) {
	// 1. Read API URL and credentials from environment variables
	whatsappApiUrl := w.cfg.WhatsappConfig.Host
	whatsappUser := w.cfg.WhatsappConfig.Username
	whatsappPass := w.cfg.WhatsappConfig.Password
	if whatsappApiUrl == "" || whatsappUser == "" || whatsappPass == "" {
		return fmt.Errorf("required WhatsApp API environment variables are not set")
	}

	// 2. Prepare the request payload
	payload := map[string]interface{}{ // Use map[string]interface{} to handle various types
		"phone":   phoneNumber,
		"message": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 3. Create a new HTTP request
	req, err := http.NewRequest("POST", whatsappApiUrl+"/send/message", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 4. Set the required headers
	req.Header.Set("Content-Type", "application/json")

	// 5. Add the Basic Authentication header
	authString := fmt.Sprintf("%s:%s", whatsappUser, whatsappPass)
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Set("Authorization", authHeader)

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message to WhatsApp API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WhatsApp API returned non-OK status: %s", resp.Status)
	}

	fmt.Println("WhatsApp message sent successfully.")
	return
}
