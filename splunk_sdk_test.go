package splunksdk_go

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetMetric(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// params for the request
	params := RequestParams{
		SearchQuery: "source=/opt/splunk/var/log/secure.log sourcetype=osx_secure |stats count",
	}
	jsonResponsePOST := `{
		"sid": "10"
	}`

	jsonResponseGET := `{
		"results":[{"count":"1250"}]
	}`

	responses := make([]map[string]interface{}, 2)
	responses[0] = map[string]interface{}{
		"POST": jsonResponsePOST,
	}
	responses[1] = map[string]interface{}{
		"GET": jsonResponseGET,
	}
	server := MutitpleMockRequest(responses)
	defer server.Close()

	spReq := SplunkRequest{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(1) * time.Second,
		},
		Params: params,
	}
	sc := SplunkCreds{
		Host:     strings.Split(strings.Split(server.URL, ":")[1], "//")[1],
		Port:     strings.Split(server.URL, ":")[2],
		Token:    "apiToken",
		Endpoint: "",
	}

	metric, err := GetMetricFromNewJob(&spReq, &sc)
	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}

	expectedMetric := 1250
	if metric != float64(expectedMetric) {
		t.Errorf("Expected %v but got %v.", expectedMetric, metric)
	}
}
func TestCreateJob(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// params for the request
	params := RequestParams{
		SearchQuery: "source=/opt/splunk/var/log/secure.log sourcetype=osx_secure |stats count",
	}
	jsonResponsePOST := `{
		"sid": "10"
	}`
	server := MockRequest(jsonResponsePOST)
	defer server.Close()

	spReq := SplunkRequest{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(1) * time.Second,
		},
		Params: params,
		// Headers: make(map[string]string),
	}
	sc := SplunkCreds{
		Host:     strings.Split(strings.Split(server.URL, ":")[1], "//")[1],
		Port:     strings.Split(server.URL, ":")[2],
		Token:    "apiToken",
		Endpoint: "",
	}

	endpoint, err := CreateJobEndpoint(&sc)
	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}
	sc.Endpoint = endpoint
	sid, err := CreateJob(&spReq, &sc)

	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}

	expectedSID := "10"
	if sid != expectedSID {
		t.Errorf("Expected %v but got %v.", expectedSID, sid)
	}
}

func TestRetrieveJobResult(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// params for the request
	params := RequestParams{}

	jsonResponseGET := `{
		"results":[{"count":"1250"}]
	}`
	server := MockRequest(jsonResponseGET)
	defer server.Close()

	spReq := SplunkRequest{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(1) * time.Second,
		},
		Params: params,
	}
	sc := SplunkCreds{
		Host:     strings.Split(strings.Split(server.URL, ":")[1], "//")[1],
		Port:     strings.Split(server.URL, ":")[2],
		Token:    "apiToken",
		Endpoint: "",
	}

	endpoint, err := CreateJobEndpoint(&sc)
	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}
	sc.Endpoint = endpoint
	sid, err := RetrieveJobResult(&spReq, &sc)

	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}

	expectedRes := make([]map[string]string, 1)
	expectedRes[0] = map[string]string{
		"count": "1250",
	}

	if sid[0]["count"] != expectedRes[0]["count"] {
		t.Errorf("Expected %v but got %v.", expectedRes, sid)
	}
}

func MockRequest(response string) *httptest.Server {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(response))
	}))
	return server
}
func MutitpleMockRequest(responses []map[string]interface{}) *httptest.Server {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
