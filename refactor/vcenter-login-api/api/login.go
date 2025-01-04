package api

import (
	"fmt"
	"io"
	"net/http"
)

func doLoginRequest(client http.Client, url string, username string, password string) (string, error) {
	// Membuat request HTTP
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("error membuat request HTTP: %w", err)
	}

	// Mengatur Basic Auth
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error mengirimkan request HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Membaca response HTTP
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error membaca body respons: %w", err)
	}

	// Log status dan body respons
	fmt.Printf("Login Response Status: %s\n", resp.Status)
	//fmt.Printf("Login Response Body: %s\n", body)

	// Mengembalikan session ID dari body respons
	sessionID := string(body) // Mengambil body sebagai session ID
	return sessionID, nil
}
