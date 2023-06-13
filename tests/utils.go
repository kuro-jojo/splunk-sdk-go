package tests

import (
	"net/http"
	"net/http/httptest"
)

func MockRequest(response string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(response))
	}))
	return server
}

func MutitpleMockRequest(responses []map[string]interface{}) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		for _, response := range responses {
			if response["POST"] != nil && r.Method == "POST" {
				_, _ = w.Write([]byte(response["POST"].(string)))
			}
			if response["GET"] != nil && r.Method == "GET" {
				_, _ = w.Write([]byte(response["GET"].(string)))
			}
		}
	}))
	return server
}
