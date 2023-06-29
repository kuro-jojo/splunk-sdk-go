package jobs

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	splunk "github.com/kuro-jojo/splunk-sdk-go/client"
)

const PATH_JOBS_V2 = "services/search/v2/jobs/"
const PATH_SAVED_SEARCHES = "services/saved/searches/"

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

func CreateServiceEndpoint(client *splunk.SplunkClient, service string) {
	host := client.Host
	port := client.Port

	if strings.HasPrefix(host, "https://") {
		host = strings.Replace(host, "https://", "", 1)
	} else if strings.HasPrefix(host, "http://") {
		host = strings.Replace(host, "http://", "", 1)
	}
	client.Endpoint = "https://" + net.JoinHostPort(host, port) + "/" + service
	client.Endpoint = strings.ReplaceAll(client.Endpoint, " ", "")
}

func PostJob(client *splunk.SplunkClient, spRequest *splunk.SplunkRequest) (*http.Response, error) {

	return HttpSearchRequest(client, "POST", spRequest)
}

func PostAlert(client *splunk.SplunkClient, spAlert *splunk.SplunkAlert) (*http.Response, error) {

	return HttpAlertRequest(client, "POST", spAlert)
}

func GetJob(client *splunk.SplunkClient) (*http.Response, error) {

	return HttpSearchRequest(client, "GET", nil)
}

func HttpSearchRequest(client *splunk.SplunkClient, method string, spRequest *splunk.SplunkRequest) (*http.Response, error) {

	if spRequest == nil {
		spRequest = &splunk.SplunkRequest{}
	}

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

func HttpAlertRequest(client *splunk.SplunkClient, method string, spAlert *splunk.SplunkAlert) (*http.Response, error) {

	if spAlert == nil {
		spAlert = &splunk.SplunkAlert{}
	}

	spAlert.Params.OutputMode = "json"

	// parameters of the request
	params := url.Values{}
	params.Add("output_mode", spAlert.Params.OutputMode)

	if method == "POST" {
		params.Add("name", spAlert.Params.Name)
		params.Add("search", spAlert.Params.SearchQuery)
		params.Add("cron_schedule", spAlert.Params.CronSchedule)
		params.Add("alert_condition", spAlert.Params.AlertCondition)
		params.Add("actions", spAlert.Params.Actions)
		params.Add("action.webhook.param.url", spAlert.Params.WebhookUrl)

		
		if spAlert.Params.EarliestTime != "" {
			params.Add("description", spAlert.Params.Description)
		}
		if spAlert.Params.EarliestTime != "" {
			params.Add("earliest_time", spAlert.Params.EarliestTime)
		}
		if spAlert.Params.LatestTime != "" {
			params.Add("latest_time", spAlert.Params.LatestTime)
		}
	}

	return splunk.MakeHttpRequest(client, method, spAlert, params)
}
