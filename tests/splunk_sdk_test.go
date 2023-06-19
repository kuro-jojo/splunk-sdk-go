package tests

import (
	"crypto/tls"
	"net/http"
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
		SearchQuery:  "source=\"http:podtato-error\" (index=\"keptn-splunk-dev\") \"[error]\" earliest=\"6/14/2023:18:00:00\" latest=\"6/15/2023:8:00:00\" | stats count",
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
	server := MutitpleMockRequest(responses, true)
	defer server.Close()

	spReq := splunk.SplunkRequest{
		Params: params,
	}
	client := splunk.SplunkClient{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(60) * time.Second,
		},
		Host:     "172.29.226.241",
		Port:     "8089",
		Token:    "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjoiSFM1MTIiLCJ2ZXIiOiJ2MiIsInR0eXAiOiJzdGF0aWMifQ.eyJpc3MiOiJhZG1pbiBmcm9tIE5DRUwxNDExOTIiLCJzdWIiOiJhZG1pbiIsImF1ZCI6ImtlcHRuIiwiaWRwIjoiU3BsdW5rIiwianRpIjoiODBkOGFkNDQ4MWY3NWQwOTYzMjY3ZWM3NjAzNjQ1NDg4NDI0ZWE1YTkyZDk0NTYzNGRkNTk1NzU1YTk3YzEyZCIsImlhdCI6MTY4NTYwNTM2MywiZXhwIjoxNjg4MTk3MzYzLCJuYnIiOjE2ODU2MDUzNjN9.eLqkWeU6TQzmfMwoJY3E0USL36pxzUri7mst-HrQb2Ay3UgZpCBfUdEM6BZ-Qgfm1gLxvGWKBsqDPGezBeiuhg",
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
	server := MockRequest(jsonResponsePOST, true)
	defer server.Close()

	spReq := splunk.SplunkRequest{
		Params: params,
	}
	client := splunk.SplunkClient{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(60) * time.Second,
		},
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
	server := MockRequest(jsonResponseGET, true)
	defer server.Close()

	spReq := splunk.SplunkRequest{
		Params: params,
	}
	client := splunk.SplunkClient{
		Client: &http.Client{
			Transport: tr,
			Timeout:   time.Duration(60) * time.Second,
		},
		Host:       strings.Split(strings.Split(server.URL, ":")[1], "//")[1],
		Port:       strings.Split(server.URL, ":")[2],
		SessionKey: "sessionKey",
		Endpoint:   "",
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
