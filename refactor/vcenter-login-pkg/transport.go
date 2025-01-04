package main

import (
	"fmt"
	"net/http"
	"net/http/httputil" // Impor package httputil untuk DumpRequest
	"strings"
)

type MyJWTTransport struct {
	transport http.RoundTripper
	sessionId string
}

func (m *MyJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.sessionId != "" {
		m.sessionId = strings.Trim(m.sessionId, "\"")
		req.Header.Set("vmware-api-session-id", m.sessionId)
		req.Header.Set("Content-Type", "application/json")
	}

	// Dump request untuk debugging
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Request Dump in MyJWTTransport:\n%s\n", dump)

	return m.transport.RoundTrip(req)
}
