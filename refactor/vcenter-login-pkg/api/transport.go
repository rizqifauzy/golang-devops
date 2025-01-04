package api

import (
	"fmt"
	"net/http"
	"net/http/httputil" // Impor package httputil untuk DumpRequest
	"strings"
)

type MyJWTTransport struct {
	Transport http.RoundTripper
	SessionId string
}

func (m *MyJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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
