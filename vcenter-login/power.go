package main

import (
	"crypto/tls" // Import package tls
	"fmt"
	"io"
	"net/http"
	"net/http/httputil" // Impor package httputil untuk DumpRequest
	"strings"
)

// Fungsi DoPowerRequest untuk menghidupkan atau mematikan VM
func DoPowerRequest(url string, sessionId string, action string) (int, error) {
	// Membuat URL dengan parameter action
	powerUrl := fmt.Sprintf("%s?action=%s", url, action)

	// Cetak URL yang akan digunakan
	fmt.Printf("Power URL: %s\n", powerUrl)

	// Membuat body kosong karena tidak diperlukan
	body := strings.NewReader("") // Body kosong

	// Membuat request HTTP
	req, err := http.NewRequest("POST", powerUrl, body)
	if err != nil {
		return 0, err
	}

	// Menambahkan header vmware-api-session-id
	sessionId = strings.Trim(sessionId, "\"")
	req.Header.Set("vmware-api-session-id", sessionId)
	req.Header.Set("Content-Type", "application/json")

	// Cetak header untuk debugging
	//fmt.Printf("Request Headers: %v\n", req.Header)

	// Dump request untuk debugging
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Request Dump:\n%s\n", dump)

	// Mengirimkan request HTTP dengan client yang mengabaikan verifikasi TLS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Mengabaikan verifikasi TLS
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode, err
	}
	defer resp.Body.Close()

	// Membaca response HTTP
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// Log status dan body respons
	//fmt.Printf("Power Response Status: %s\n", resp.Status)
	//fmt.Printf("Power Response Body: %s\n", bodyBytes)

	// Cek apakah ada error dalam respons
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status: ", resp.Status)
		return resp.StatusCode, fmt.Errorf("failed to power on VM: %s", bodyBytes)
	}
	fmt.Println("Status: ", resp.Status)

	return resp.StatusCode, nil
}
