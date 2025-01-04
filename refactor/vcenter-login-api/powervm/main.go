package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/rizqifauzy/golang-devops/refactor/vcenter-login-api/api"
)

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

	// Mengirimkan request HTTP dengan client yang mengabaikan verifikasi TL

	apiInstance := api.New(api.Options{
		Password: password,
		Username: username,
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/api/session",
	})

	vmUrl := parsedURL.String() + fmt.Sprintf("/api/vcenter/vm/%s/power", vmName)

	message, err := apiInstance.DoPowerRequest(vmUrl, "stop")

	if err != nil {
		if requestErr, ok := err.(api.RequestError); ok {
			fmt.Printf("%s (HTTP Code: %d, Body: %s)\n", requestErr.Error(), requestErr.HTTPCode, requestErr.Body)
			os.Exit(0)
		}
		log.Fatal(err)
	}

	fmt.Printf("Response: \n%s\n", message.GetResponse())
}
