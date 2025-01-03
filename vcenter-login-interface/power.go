package main

import (
	"crypto/tls" // Import package tls
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil" // Impor package httputil untuk DumpRequest
	"strings"
)

// Fungsi DoPowerRequest untuk menghidupkan atau mematikan VM
func DoPowerRequest(url string, sessionId string, action string) (Output, error) {
	// Membuat URL dengan parameter action
	powerUrl := fmt.Sprintf("%s?action=%s", url, action)

	// Cetak URL yang akan digunakan
	fmt.Printf("Power URL: %s\n", powerUrl)

	// Membuat body kosong karena tidak diperlukan
	body := strings.NewReader("") // Body kosong

	// Membuat request HTTP
	req, err := http.NewRequest("POST", powerUrl, body)
	if err != nil {
		return nil, err
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
		return nil, err
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
		//return nil, fmt.Errorf("Invalid: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Membaca response HTTP
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Cek apakah ada error dalam respons
	if resp.StatusCode == 204 {
		// Jika status tidak OK, kembalikan respons JSON sebagai string
		jsonResponse := JSONResponse{
			JSON: string(bodyBytes),
		}
		return jsonResponse, nil
	}

	var response Response

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %s", err)
	}

	// Log status dan body respons
	//fmt.Printf("Power Response Status: %s\n", resp.Status)
	//fmt.Printf("Power Response Body: %s\n", bodyBytes)

	// Cek apakah ada error dalam respons
	if resp.StatusCode == 400 {
		fmt.Println("Status: ", resp.Status)
		//return nil, fmt.Errorf("failed to power on VM: %s", bodyBytes)
		return response, nil
	}
	//fmt.Println("Status: ", resp.Status)

	return response, nil
}
