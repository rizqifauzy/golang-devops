package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// password dikirm dalam format json, dan token dikembalikan juga dalam format json
type LoginRequest struct {
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(loginURL, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", err)
	}

	response, err := http.Post(loginURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", fmt.Errorf("http Post error: %s", err)
	}

	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Invalid output (HTTP Code %d): %s\n", response.StatusCode, string(resBody))
	}

	if !json.Valid(resBody) {
		return "", RequestError{
			Err:      fmt.Sprintf("Response is not a json"),
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
		}
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", RequestError{
			Err:      fmt.Sprintf("Page unmarshal error: %s", err),
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
		}
	}

	if loginResponse.Token == "" {
		return "", RequestError{
			Err:      fmt.Sprintf("Token is empty"),
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
		}
	}

	return loginResponse.Token, nil
}
