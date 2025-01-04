package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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

// Definisikan struct baru untuk mengimplementasikan antarmuka Output
// type JSONResponse struct {
// 	JSON string `json:"json"`
// }

// // Implementasikan metode GetResponse pada struct JSONResponse
// func (jr JSONResponse) GetResponse() string {
// 	//return jr.JSON
// 	return "Kode berhasil dijalankan"
// }

func main() {
	var (
		vmName string
		//parsedURL *url.URL
		//err       error
	)

	vcenterUrl := os.Getenv("VSPHERE_SERVER")
	username := os.Getenv("VSPHERE_USER")
	password := os.Getenv("VSPHERE_PASSWORD")

	if vcenterUrl == "" || username == "" || password == "" {
		log.Fatal("VSPHERE_SERVER, VSPHERE_USER, atau VSPHERE_PASSWORD tidak diatur")
	}

	flag.StringVar(&vmName, "vm", "", "Nama VM")
	flag.Parse()

	if vmName == "" {
		log.Fatal("Nama VM tidak diberikan")
	}

	parsedURL, err := url.ParseRequestURI(vcenterUrl)
	if err != nil {
		fmt.Printf("Help: ./http-get -h\nURL is not valid URL: %s\n", vcenterUrl)
		os.Exit(1)
	}

	// Mengirimkan request HTTP dengan client yang mengabaikan verifikasi TLS
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{
		Transport: transport,
	}

	sessionId, err := DoLoginRequest(client, parsedURL.Scheme+"://"+parsedURL.Host+"/api/session", username, password)
	if err != nil {
		log.Fatal(err)
	}

	if sessionId == "" {
		log.Fatal("Gagal mendapatkan session ID")
	}

	client.Transport = &MyJWTTransport{
		transport: transport,
		sessionId: sessionId,
	}

	vmUrl := parsedURL.String() + fmt.Sprintf("/api/vcenter/vm/%s/power", vmName)

	message, err := DoPowerRequest(client, vmUrl, sessionId, "stop") // Ganti "stop" dengan "start"

	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
			fmt.Printf("%s (HTTP Code: %d, Body: %s)\n", requestErr.Error(), requestErr.HTTPCode, requestErr.Body)
			os.Exit(0)
		}
		log.Fatal(err)
	}

	fmt.Printf("Response: \n%s\n", message.GetResponse())
}
