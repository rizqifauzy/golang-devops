package api

import "net/http"

type Options struct {
	Password string
	LoginURL string
}

type APIIface interface {
	DoGetRequest(requestURL string) (Response, error)
}

type API struct {
	Options Options
	Client  http.Client
}

func New(options Options) APIIface {
	return API{
		Options: options,
		Client: http.Client{
			Transport: MyJWTTransport{
				transport: http.DefaultTransport,
				password:  options.Password,
				loginURL:  options.LoginURL,
			},
		},
	}
}
