package tests

import (
	"net/http"
	"net/http/httptest"
)

// mock an http server 
func MockRequest(response string, verify bool) *httptest.Server {
	server := &httptest.Server{}
	if verify {
		server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			writeResponses(response, &w, r)
		}))

	} else {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			writeResponses(response, &w, r)
		}))
	}
	return server
}

func MutitpleMockRequest(responses []map[string]interface{}, verify bool) *httptest.Server {
	server := &httptest.Server{}
	if verify {
		server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			writeResponses(responses, &w, r)
		}))

	} else {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			writeResponses(responses, &w, r)
		}))
	}
	return server
}

func writeResponses(responses interface{}, w *http.ResponseWriter, r *http.Request) {

	switch responses.(type) {
	case []map[string]interface{}:
		for _, response := range responses.([]map[string]interface{}) {
			if response["POST"] != nil && r.Method == "POST" {
				_, _ = (*w).Write([]byte(response["POST"].(string)))
			}
			if response["GET"] != nil && r.Method == "GET" {
				_, _ = (*w).Write([]byte(response["GET"].(string)))
			}
		}
	case string:
		_, _ = (*w).Write([]byte(responses.(string)))
	}
}
