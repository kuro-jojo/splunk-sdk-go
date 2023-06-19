package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	splunk "github.com/kuro-jojo/splunk-sdk-go/client"
	job "github.com/kuro-jojo/splunk-sdk-go/jobs"
)

func TestGetMetric(t *testing.T) {
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

	client := splunk.NewClientAuthenticatedByToken(
		&http.Client{
			Timeout: time.Duration(60) * time.Second,
		},
		getServerHostname(server),
		getServerPort(server),
		"token",
		true,
	)

	defer server.Close()

	spReq := splunk.SplunkRequest{
		Params: splunk.RequestParams{
			SearchQuery: "source=\"http:podtato-error\" (index=\"keptn-splunk-dev\") \"[error]\" earliest=\"2023-06-15T15:04:45.000Z\" latest=-3d | stats count",
		},
	}

	metric, err := job.GetMetricFromNewJob(client, &spReq)

	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}

	expectedMetric := 1250
	if metric != float64(expectedMetric) {
		t.Errorf("Expected %v but got %v.", expectedMetric, metric)
	}
}

func TestCreateJob(t *testing.T) {

	jsonResponsePOST := `{
		"sid": "10"
	}`
	server := MockRequest(jsonResponsePOST, true)
	defer server.Close()

	spReq := splunk.SplunkRequest{
		Params: splunk.RequestParams{
			SearchQuery: "source=\"http:podtato-error\" (index=\"keptn-splunk-dev\") \"[error]\" earliest=\"2023-06-15T15:04:45.000Z\" latest=-3d | stats count",
		},
	}
	client := splunk.NewClientAuthenticatedByToken(
		&http.Client{
			Timeout: time.Duration(60) * time.Second,
		},
		getServerHostname(server),
		getServerPort(server),
		"token",
		true,
	)

	endpoint, err := job.CreateJobEndpoint(client)
	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}
	client.Endpoint = endpoint
	sid, err := job.CreateJob(client, &spReq)

	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}

	expectedSID := "10"
	if sid != expectedSID {
		t.Errorf("Expected %v but got %v.", expectedSID, sid)
	}
}

func TestRetrieveJobResult(t *testing.T) {

	jsonResponseGET := `{
		"results":[{"count":"1250"}]
	}`
	server := MockRequest(jsonResponseGET, true)
	defer server.Close()

	spReq := splunk.SplunkRequest{
		Params: splunk.RequestParams{
			SearchQuery: "source=\"http:podtato-error\" (index=\"keptn-splunk-dev\") \"[error]\" earliest=\"2023-06-15T15:04:45.000Z\" latest=-3d | stats count",
		},
	}
	client := splunk.NewClientAuthenticatedByToken(
		&http.Client{
			Timeout: time.Duration(60) * time.Second,
		},
		getServerHostname(server),
		getServerPort(server),
		"token",
		true,
	)

	endpoint, err := job.CreateJobEndpoint(client)
	if err != nil {
		t.Fatalf("Got an error : %s", err)
	}
	client.Endpoint = endpoint
	sid, err := job.RetrieveJobResult(client, &spReq)

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

func getServerHostname(server *httptest.Server) string {
	host := strings.Split(strings.Split(server.URL, ":")[1], "//")[1]

	return host
}

func getServerPort(server *httptest.Server) string {
	port := strings.Split(server.URL, ":")[2]

	return port
}
