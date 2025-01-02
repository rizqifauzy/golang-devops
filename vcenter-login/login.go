package main

import (
	"crypto/tls" // Import package tls
	"fmt"        // Import fmt untuk logging
	"io"
	"net/http"
)

// Fungsi DoLoginRequest untuk melakukan login ke vCenter
func DoLoginRequest(url string, username string, password string) (string, error) {
	// Membuat request HTTP
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	// Mengatur Basic Auth
	req.SetBasicAuth(username, password)

	// Mengirimkan request HTTP dengan client yang mengabaikan verifikasi TLS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Mengabaikan verifikasi TLS
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Membaca response HTTP
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Log status dan body respons
	fmt.Printf("Login Response Status: %s\n", resp.Status)
	//fmt.Printf("Login Response Body: %s\n", body)

	// Mengembalikan session ID dari body respons
	sessionID := string(body) // Mengambil body sebagai session ID
	return sessionID, nil
}
