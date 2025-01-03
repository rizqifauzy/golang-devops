package main

import (
	"encoding/json"
	"fmt"
	"log"
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

// func (r Response) GetResponse() string {
// 	words := []string{}
// 	for _, msg := range r.Messages {
// 		words = append(words, msg.DefaultMessage)
// 	}
// 	return fmt.Sprintf("Default Messages: %s", strings.Join(words, ", "))
// }

func (r Response) GetResponse() string {
	if len(r.Messages) == 0 {
		return "No messages available"
	}
	return fmt.Sprintf("Status: %s\nMessage: %s", r.ErrorType, r.Messages[0].DefaultMessage)
}

func main() {
	// JSON respons yang akan diurai
	jsonData := `{
		"error_type": "ALREADY_IN_DESIRED_STATE",
		"messages": [
			{
				"args": [],
				"default_message": "Virtual machine is already powered off.",
				"id": "com.vmware.api.vcenter.vm.power.already_powered_off"
			},
			{
				"args": [],
				"default_message": "The attempted operation cannot be performed in the current state (Powered off).",
				"id": "vmsg.InvalidPowerState.summary"
			}
		]
	}`

	res, err := getJson(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: \n%s\n", res.GetResponse())

}

func getJson(jsonData string) (Output, error) {
	var response Response

	// Unmarshal JSON ke dalam struct Response
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %v", err)
	}

	return response, nil
}
