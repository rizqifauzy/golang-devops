package api

import (
	"crypto/tls"
	"log"
	"net/http"
)

type Options struct {
	Username string
	Password string
	LoginURL string
}

type APIIface interface {
	DoPowerRequest(requestURL string, action string) (Output, error)
}

type api struct {
	Options Options
	Client  http.Client
}

func New(options Options) APIIface {
	if options.LoginURL == "" {
		log.Fatal("LoginURL is required")
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return api{
		Options: options,
		Client: http.Client{
			Transport: &MyJWTTransport{
				Transport: transport,
				loginURL:  options.LoginURL,
				username:  options.Username,
				password:  options.Password,
			},
		},
	}
}
