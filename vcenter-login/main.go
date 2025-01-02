package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var vmName string
	flag.StringVar(&vmName, "vm", "", "Nama VM")
	flag.Parse()

	if vmName == "" {
		log.Fatal("Nama VM tidak diberikan")
	}

	vcenterUrl := os.Getenv("VSPHERE_SERVER")
	username := os.Getenv("VSPHERE_USER")
	password := os.Getenv("VSPHERE_PASSWORD")

	if vcenterUrl == "" || username == "" || password == "" {
		log.Fatal("VSPHERE_SERVER, VSPHERE_USER, atau VSPHERE_PASSWORD tidak diatur")
	}

	sessionId, err := DoLoginRequest(vcenterUrl+"/api/session", username, password)
	if err != nil {
		log.Fatal(err)
	}

	if sessionId == "" {
		log.Fatal("Gagal mendapatkan session ID")
	}

	vmUrl := vcenterUrl + fmt.Sprintf("/api/vcenter/vm/%s/power", vmName)

	statusCode, err := DoPowerRequest(vmUrl, sessionId, "stop") // Ganti "stop" dengan "start"
	//if err != nil {
	//	log.Fatal(err)
	//}
	if statusCode == 0 {
		fmt.Println("Terjadi kesalahan saat menjalankan kode")

	}
	if statusCode == 400 {
		fmt.Println("VM sudah dalam kondisi DESIRED STATE")

	}
	if statusCode == 204 {
		fmt.Println("Kode berhasil dijalankan")

	}
	//fmt.Println("Status Code: ", statusCode)
	//fmt.Println("Kode berhasil dijalankan")
}
