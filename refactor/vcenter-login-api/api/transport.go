package api

import (
	"crypto/tls"
	"fmt"
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
	// Jika session ID belum ada, lakukan login
	if m.SessionId == "" {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		sessionId, err := doLoginRequest(*client, m.loginURL, m.username, m.password)
		if err != nil {
			return nil, fmt.Errorf("error login: %w", err)
		}

		if sessionId == "" {
			return nil, fmt.Errorf("gagal mendapatkan session ID")
		}

		m.SessionId = strings.Trim(sessionId, "\"")
		fmt.Printf("Session ID berhasil didapatkan: %s\n", m.SessionId)
	}

	if m.SessionId != "" {
		// Tambahkan session ID ke header
		req.Header.Set("vmware-api-session-id", m.SessionId)
		req.Header.Set("Content-Type", "application/json")
	}

	// Dump request untuk debugging
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Request Dump in MyJWTTransport:\n%s\n", dump)

	// Lakukan request
	return m.Transport.RoundTrip(req)
}
