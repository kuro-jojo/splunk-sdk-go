package jobs

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	splunk "github.com/kuro-jojo/splunk-sdk-go/client"
)

const PATH_JOBS_V2 = "services/search/v2/jobs/"

func ValidateSearchQuery(searchQuery string) string {
	// the search must start with the "search" keyword
	const QUERY_PREFIX = "search "
	if !strings.HasPrefix(searchQuery, QUERY_PREFIX) {
		return QUERY_PREFIX + searchQuery
	}
	return searchQuery
}

func HttpJobRequest(client *splunk.SplunkClient, method string, spRequest *splunk.SplunkRequest) (*http.Response, error) {

	if spRequest.Params.OutputMode == "" {
		spRequest.Params.OutputMode = "json"
	}
	if spRequest.Params.ExecMode == "" {
		spRequest.Params.ExecMode = "blocking"
	}

	// parameters of the request
	params := url.Values{}
	params.Add("output_mode", spRequest.Params.OutputMode)
	params.Add("exec_mode", spRequest.Params.ExecMode)

	if method == "POST" {
		params.Add("search", spRequest.Params.SearchQuery)
		if spRequest.Params.EarliestTime != "" {
			params.Add("earliest_time", spRequest.Params.EarliestTime)
		}
		if spRequest.Params.LatestTime != "" {
			params.Add("latest_time", spRequest.Params.LatestTime)
		}
	}

	return splunk.MakeHttpRequest(client, method, spRequest, params)
}

func CreateJobEndpoint(client *splunk.SplunkClient) (string, error) {
	host := client.Host
	port := client.Port

	if strings.HasPrefix(host, "https://") {
		host = strings.Replace(host, "https://", "", 1)
	} else if strings.HasPrefix(host, "http://") {
		host = strings.Replace(host, "http://", "", 1)
	}
	return "https://" + net.JoinHostPort(host, port) + "/" + PATH_JOBS_V2, nil
}

func PostJob(client *splunk.SplunkClient, spRequest *splunk.SplunkRequest) (*http.Response, error) {

	return HttpSearchRequest(client, "POST", spRequest)
}

func GetJob(client *splunk.SplunkClient, spRequest *splunk.SplunkRequest) (*http.Response, error) {

	return HttpSearchRequest(client, "GET", spRequest)
}

func HttpSearchRequest(client *splunk.SplunkClient, method string, spRequest *splunk.SplunkRequest) (*http.Response, error) {

	spRequest.Params.OutputMode = "json"
	spRequest.Params.ExecMode = "blocking"

	// parameters of the request
	params := url.Values{}
	params.Add("output_mode", spRequest.Params.OutputMode)
	params.Add("exec_mode", spRequest.Params.ExecMode)

	if method == "POST" {
		params.Add("search", spRequest.Params.SearchQuery)
		if spRequest.Params.EarliestTime != "" {
			params.Add("earliest_time", spRequest.Params.EarliestTime)
		}
		if spRequest.Params.LatestTime != "" {
			params.Add("latest_time", spRequest.Params.LatestTime)
		}
	}

	return splunk.MakeHttpRequest(client, method, spRequest, params)
}
