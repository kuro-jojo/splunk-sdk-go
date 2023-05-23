package splunksdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const CREATE_JOB_URI = "services/search/jobs/"

type SplunkCreds struct {
	Host     string
	Port     string
	Token    string
	Endpoint string
}

func GetMetric(client *http.Client, splunkCreds *SplunkCreds, searchQuery string, headers map[string]string) (float64, error) {

	const RESULTS_URI = "results"
	endpoint, err := CreateJobEndpoint(splunkCreds)
	if err != nil {
		return -1, err
	}
	splunkCreds.Endpoint = endpoint
	searchQuery = validateSearchQuery(searchQuery)
	sid, err := postSearch(client, splunkCreds, searchQuery, headers)
	if err != nil {
		return -1, fmt.Errorf("error : %s", err)
	}

	newEndpoint := endpoint + sid
	// check if the endpoint is correctly formed
	if !strings.HasSuffix(newEndpoint, "/") {
		newEndpoint += "/"
	}

	// the endpoint where to find the corresponding job
	splunkCreds.Endpoint = newEndpoint + RESULTS_URI
	fmt.Printf("Endpoint : %s\n\n", splunkCreds.Endpoint)

	res, err := getSearch(client, splunkCreds, headers)
	if err != nil {
		return -1, fmt.Errorf("error while handling the results. Error message : %s", err)
	}
	// if the result is not a metric
	if len(res) != 1 {
		return -1, fmt.Errorf("incorrect search result. Error message : %s", err)
	}
	var metrics []string

	for _, v := range res[0] {
		metrics = append(metrics, v)
	}
	metric, err := strconv.ParseFloat(metrics[0], 64)
	if err != nil {
		return -1, fmt.Errorf("convert metric to float failed. Error message : %s", err)
	}

	return metric, nil
}

func CreateJobEndpoint(sc *SplunkCreds) (string, error) {
	host := sc.Host
	port := sc.Port

	match := `^((localhost)|(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})|([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}))$`

	if !regexp.MustCompile(match).MatchString(host) {
		return "", fmt.Errorf("")
	}
	return "http://" + host + ":" + port + "/" + CREATE_JOB_URI, nil
}

func getSID(resp []byte) (string, error) {
	respJson := string(resp)

	var sid map[string]string
	json.Unmarshal([]byte(respJson), &sid)
	if len(sid) <= 0 {
		return "", fmt.Errorf("no sid found")
	}
	return sid["sid"], nil
}

func getBearer(token string) string {
	if !strings.HasPrefix(token, "Bearer") {
		return "Bearer " + token
	}
	return token
}

func httpRequest(method string, client *http.Client, splunkCreds *SplunkCreds, params url.Values, headers map[string]string) (*http.Response, error) {

	// create a new request
	req, err := http.NewRequest(method, splunkCreds.Endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	// add the headers
	for header, val := range headers {
		req.Header.Add(header, val)
	}
	if req.Header.Get("Authorization") != "" {
		req.Header.Add("Authorization", getBearer(splunkCreds.Token))
	} else {
		req.Header.Set("Authorization", getBearer(splunkCreds.Token))
	}
	// get the response
	resp, err := client.Do(req)
	
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()
	return resp, nil
}
func post(client *http.Client, splunkCreds *SplunkCreds, params url.Values, headers map[string]string) (*http.Response, error) {

	return httpRequest("POST", client, splunkCreds, params, headers)
}

func get(client *http.Client, splunkCreds *SplunkCreds, params url.Values, headers map[string]string) (*http.Response, error) {

	return httpRequest("GET", client, splunkCreds, params, headers)
}

func validateSearchQuery(searchQuery string) string {
	// the search must start with the "search" keyword
	const QUERY_PREFIX = "search "
	if !strings.HasPrefix(searchQuery, QUERY_PREFIX) {
		return QUERY_PREFIX + searchQuery
	}
	return searchQuery
}

func postSearch(client *http.Client, splunkCreds *SplunkCreds, searchQuery string, headers map[string]string) (string, error) {
	// make the post request
	// params to send
	params := url.Values{}
	params.Add("search", searchQuery)
	params.Add("output_mode", "json")
	params.Add("exec_mode", "blocking")

	resp, err := post(client, splunkCreds, params, headers)
	if err != nil {
		return "", fmt.Errorf("error : %s", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error : %s", err)
	}

	// create the new endpoint for the get request
	sid, err := getSID(body)
	if err != nil {
		return "", fmt.Errorf("error : %s", err)
	}

	return sid, nil
}

func getSearch(client *http.Client, splunkCreds *SplunkCreds, headers map[string]string) ([]map[string]string, error) {
	// new parameters for the get request
	getParams := url.Values{}
	getParams.Add("output_mode", "json")

	// make the get request
	getResp, err := get(client, splunkCreds, getParams, headers)
	if err != nil {
		return nil, fmt.Errorf("error : %s", err)
	}
	// get the body of the response
	getBody, err := io.ReadAll(getResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error : %s", err)
	}

	// get the metric we want
	var results map[string][]map[string]string
	json.Unmarshal([]byte(getBody), &results)

	return results["results"], nil
}
