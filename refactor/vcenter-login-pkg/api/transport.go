package api

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type MyJWTTransport struct {
	Transport http.RoundTripper
	SessionId string
	username  string
	password  string
	loginURL  string
}

func (m *MyJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.SessionId == "" {
		if m.password != "" {
			// Buat HTTP client dengan InsecureSkipVerify
			tlsConfig := &tls.Config{InsecureSkipVerify: true} // Skip SSL verification
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: tlsConfig,
				},
			}

			sessionId, err := doLoginRequest(*client, m.loginURL, m.username, m.password)
			if err != nil {
				log.Fatal(err)
			}

			if sessionId == "" {
				log.Fatal("Gagal mendapatkan session ID")
			}

			m.SessionId = sessionId
		}
	}
	if m.SessionId != "" {
		m.SessionId = strings.Trim(m.SessionId, "\"")
		req.Header.Set("vmware-api-session-id", m.SessionId)
		req.Header.Set("Content-Type", "application/json")
	}

	// Dump request untuk debugging
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Request Dump in MyJWTTransport:\n%s\n", dump)

	return m.Transport.RoundTrip(req)
}
