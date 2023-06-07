package splunksdk_go

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	splunk "github.com/kuro-jojo/splunk-sdk-go/client"
	job "github.com/kuro-jojo/splunk-sdk-go/jobs"
)

func TestGetMetric(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// params for the request
	params := splunk.RequestParams{
		SearchQuery:  "source=/opt/splunk/var/log/secure.log sourcetype=osx_secure |stats count",
		EarliestTime: "-5m",
		LatestTime:   "-1m",
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

	spReq := splunk.SplunkRequest{
		Params: params,
	}
	client := splunk.SplunkClient{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(60) * time.Second,
		},
		// Host:  "172.29.226.241",
		// Port:  "8089",
		// Token: "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJhZG1pbiBmcm9tIE5DRUwxNDExOTIiLCJzdWIiOiJhZG1pbiIsImF1ZCI6InRlc3QiLCJpZHAiOiJTcGx1bmsiLCJqdGkiOiI2MTE5ZjE3NmExZmEyMmZkZjA1MTM5M2JhNDJkZTA0OTczZTBlMjFkOTRmYjcyNDdmYzQwZTAzYmJhYWIwZTdhIiwiaWF0IjoxNjg1NTM2MjIzLCJleHAiOjE2ODgxMjgyMjMsIm5iciI6MTY4NTUzNjIyM30.gx_mxwT6xdKoiP2Mrh_DsHcGHyxG9RlBusAaZlLOA9n-U8J6gmWQCMkTcvrEtR6l5LdvsLZ0BW8n06bNrAIEYw",
		// // Endpoint: "",

		Host:     strings.Split(strings.Split(server.URL, ":")[1], "//")[1],
		Port:     strings.Split(server.URL, ":")[2],
		Token:    "apiToken",
		Endpoint: "",
	}

	metric, err := job.GetMetricFromNewJob(&client, &spReq)

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
	params := splunk.RequestParams{
		SearchQuery:  "source=/opt/splunk/var/log/secure.log sourcetype=osx_secure |stats count",
		EarliestTime: "-1y@w1",
		LatestTime:   "-500h",
	}
	jsonResponsePOST := `{
		"sid": "10"
	}`
	server := MockRequest(jsonResponsePOST)
	defer server.Close()

	spReq := splunk.SplunkRequest{
		Params: params,
	}
	client := splunk.SplunkClient{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(60) * time.Second,
		},
		// Host:  "172.29.226.241",
		// Port:  "8089",
		// Token: "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJhZG1pbiBmcm9tIE5DRUwxNDExOTIiLCJzdWIiOiJhZG1pbiIsImF1ZCI6InRlc3QiLCJpZHAiOiJTcGx1bmsiLCJqdGkiOiI2MTE5ZjE3NmExZmEyMmZkZjA1MTM5M2JhNDJkZTA0OTczZTBlMjFkOTRmYjcyNDdmYzQwZTAzYmJhYWIwZTdhIiwiaWF0IjoxNjg1NTM2MjIzLCJleHAiOjE2ODgxMjgyMjMsIm5iciI6MTY4NTUzNjIyM30.gx_mxwT6xdKoiP2Mrh_DsHcGHyxG9RlBusAaZlLOA9n-U8J6gmWQCMkTcvrEtR6l5LdvsLZ0BW8n06bNrAIEYw",
		// // Endpoint: "",

		Host:     strings.Split(strings.Split(server.URL, ":")[1], "//")[1],
		Port:     strings.Split(server.URL, ":")[2],
		Token:    "apiToken",
		Endpoint: "",
	}

	endpoint, err := job.CreateJobEndpoint(&client)
	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}
	client.Endpoint = endpoint
	sid, err := job.CreateJob(&client, &spReq)

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
	params := splunk.RequestParams{}

	jsonResponseGET := `{
		"results":[{"count":"1250"}]
	}`
	server := MockRequest(jsonResponseGET)
	defer server.Close()

	spReq := splunk.SplunkRequest{
		Params: params,
	}
	client := splunk.SplunkClient{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(60) * time.Second,
		},
		// Host:  "172.29.226.241",
		// Port:  "8089",
		// Token: "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJhZG1pbiBmcm9tIE5DRUwxNDExOTIiLCJzdWIiOiJhZG1pbiIsImF1ZCI6InRlc3QiLCJpZHAiOiJTcGx1bmsiLCJqdGkiOiI2MTE5ZjE3NmExZmEyMmZkZjA1MTM5M2JhNDJkZTA0OTczZTBlMjFkOTRmYjcyNDdmYzQwZTAzYmJhYWIwZTdhIiwiaWF0IjoxNjg1NTM2MjIzLCJleHAiOjE2ODgxMjgyMjMsIm5iciI6MTY4NTUzNjIyM30.gx_mxwT6xdKoiP2Mrh_DsHcGHyxG9RlBusAaZlLOA9n-U8J6gmWQCMkTcvrEtR6l5LdvsLZ0BW8n06bNrAIEYw",
		// // Endpoint: "",

		Host:     strings.Split(strings.Split(server.URL, ":")[1], "//")[1],
		Port:     strings.Split(server.URL, ":")[2],
		SessionKey:    "Bearer apiToken",
		Endpoint: "",
	}

	endpoint, err := job.CreateJobEndpoint(&client)
	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}
	client.Endpoint = endpoint
	sid, err := job.RetrieveJobResult(&client, &spReq)

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
