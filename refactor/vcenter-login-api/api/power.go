package api

import (
	// Import package tls
	"encoding/json"
	"fmt"
	"io"
	"net/http" // Impor package httputil untuk DumpRequest
	"strings"
)

type Output interface {
	GetResponse() string
}

// Definisikan struct untuk menampung seluruh respons JSON
type Response struct {
	ErrorType string    `json:"error_type"`
	Messages  []Message `json:"messages"`
}

// Definisikan struct untuk menampung setiap pesan dalam array messages
type Message struct {
	Args           []interface{} `json:"args"`
	DefaultMessage string        `json:"default_message"`
	ID             string        `json:"id"`
}

func (r Response) GetResponse() string {
	if len(r.Messages) == 0 {
		return "No messages available"
	}
	return fmt.Sprintf("Status: %s\nMessage: %s", r.ErrorType, r.Messages[0].DefaultMessage)
}

// Fungsi DoPowerRequest untuk menghidupkan atau mematikan VM
func (a api) DoPowerRequest(url string, action string) (Output, error) {
	// Membuat URL dengan parameter action
	powerUrl := fmt.Sprintf("%s?action=%s", url, action)

	// Membuat body kosong karena tidak diperlukan
	body := strings.NewReader("") // Body kosong

	// Membuat request HTTP
	req, err := http.NewRequest("POST", powerUrl, body)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		//return nil, fmt.Errorf("Invalid: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Membaca response HTTP
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		//return nil, err
		return nil, RequestError{
			Err:      fmt.Sprintf("Error read body: %s", err),
			HTTPCode: resp.StatusCode,
			Body:     string(bodyBytes),
		}
	}

	// Cek apakah ada error dalam respons
	// if resp.StatusCode == 204 {
	// 	// Jika status tidak OK, kembalikan respons JSON sebagai string
	// 	jsonResponse := JSONResponse{
	// 		JSON: string(bodyBytes),
	// 	}
	// 	return jsonResponse, nil
	// }

	if !json.Valid(bodyBytes) {
		return nil, RequestError{
			Err:      fmt.Sprintf("Kode berhasil dijalankan"),
			HTTPCode: resp.StatusCode,
			Body:     string(bodyBytes),
		}
	}

	var response Response

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		//return nil, fmt.Errorf("unmarshal error: %s", err)
		return nil, RequestError{
			Err:      fmt.Sprintf("unmarshal error: %s", err),
			HTTPCode: resp.StatusCode,
			Body:     string(bodyBytes),
		}
	}

	// Cek apakah ada error dalam respons
	if resp.StatusCode == 400 {

		return response, nil
	}

	return response, nil
}
