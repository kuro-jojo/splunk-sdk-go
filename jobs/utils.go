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
const PATH_TRIGGERED_ALERTS = "services/alerts/fired_alerts/"

type splunkAlertEntry struct {
	Name string  `json:"name"`
}

type splunkAlertList struct {
	Item []splunkAlertEntry  `json:"entry"`
}

type TriggeredAlerts struct {
	Origin string `json:"origin"`
	Updated string `json:"updated"`
	Entry []EntryItem `json:"entry"`
}

type TriggeredInstances struct {
	Origin string `json:"origin"`
	Updated string `json:"updated"`
	Entry []EntryItem `json:"entry"`
}	

type EntryItem struct {
	Name string `json:"name"`
	Links Links `json:"links"`
	Content Content `content:"content"`
}

type Links struct {
	Alternate string `json:"alternate"`
	List string `json:"list"`
	Remove string `json:"remove"`
	Job string `json:"job"`
	SavedSearch string `json:"savedsearch"`
}

type Content struct {
	Sid string `json:"sid"`
	SavedSearchName string `json:"savedsearch_name"`
	TriggerTime int `json:"trigger_time"`
	TriggeredAlertCount int `json:"triggered_alert_count"`
}

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

func GetAlerts(client *splunk.SplunkClient) (*http.Response, error) {

	return HttpAlertRequest(client, "GET", nil)
}

func DeleteAlert(client *splunk.SplunkClient, spAlert *splunk.SplunkAlert) (*http.Response, error) {

	return HttpAlertRequest(client, "DELETE", spAlert)
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
		
		if spAlert.Params.Name != "" {
			params.Add("name", spAlert.Params.Name)
		}
		if spAlert.Params.Actions != "" {
			params.Add("actions", spAlert.Params.Actions)
		}
		if spAlert.Params.WebhookUrl != "" {
			params.Add("action.webhook.param.url", spAlert.Params.WebhookUrl)
		}
		if spAlert.Params.SearchQuery != "" {
			params.Add("search", spAlert.Params.SearchQuery)
		}
		if spAlert.Params.CronSchedule != "" {
			params.Add("cron_schedule", spAlert.Params.CronSchedule)
		}
		if spAlert.Params.AlertCondition != "" {
			params.Add("alert_condition", spAlert.Params.AlertCondition)
		}

		params.Add("is_scheduled", "1")

		if spAlert.Params.EarliestTime != "" {
			params.Add("dispatch.earliest_time", spAlert.Params.EarliestTime)
		}
		if spAlert.Params.LatestTime != "" {
			params.Add("dispatch.latest_time", spAlert.Params.LatestTime)
		}

		params.Add("alert_type", "custom")

		if spAlert.Params.Description != "" {
			params.Add("description", spAlert.Params.Description)
		}

		params.Add("alert.track", "1")
		
	}

	return splunk.MakeAlertHttpRequest(client, method, spAlert, params)
}
