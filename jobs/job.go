package jobs

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	splunk "github.com/kuro-jojo/splunk-sdk-go/client"
)

// Return a metric from a new created job
func GetMetricFromNewJob(client *splunk.SplunkClient, spRequest *splunk.SplunkRequest) (float64, error) {

	const RESULTS_URI = "results"

	endpoint, err := CreateJobEndpoint(client)
	if err != nil {
		return -1, fmt.Errorf("error while creating the endpoint : %s", err)
	}

	client.Endpoint = endpoint
	spRequest.Params.SearchQuery = ValidateSearchQuery(spRequest.Params.SearchQuery)
	sid, err := CreateJob(client, spRequest)
	if err != nil {
		return -1, fmt.Errorf("error while creating the job : %s", err)
	}

	newEndpoint := endpoint + sid
	// check if the endpoint is correctly formed
	if !strings.HasSuffix(newEndpoint, "/") {
		newEndpoint += "/"
	}

	// the endpoint where to find the corresponding job
	client.Endpoint = newEndpoint + RESULTS_URI

	res, err := RetrieveJobResult(client, spRequest)

	if err != nil {
		return -1, fmt.Errorf("error while handling the results. Error message : %s", err)
	}
	// if the result is not a metric
	if len(res) != 1 {
		return -1, fmt.Errorf("result is not a metric. Error message : %v", err)
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

// this function create a new job and return its SID
func CreateJob(client *splunk.SplunkClient, spRequest *splunk.SplunkRequest) (string, error) {

	resp, err := PostJob(client, spRequest)

	if err != nil {
		return "", fmt.Errorf("error while making the post request : %s", err)
	}

	body, err := io.ReadAll(resp.Body)
	// handle error
	if !strings.HasPrefix(strconv.Itoa(resp.StatusCode), "2") {
		status, err := splunk.HandleHttpError(body)
		if err == nil {
			return "", fmt.Errorf("http error :  %s and %v", status, spRequest.Params)
		} else {
			return "", fmt.Errorf("http error :  %s and %v", resp.Status, spRequest.Params)
		}
	}

	if err != nil {
		return "", fmt.Errorf("error while getting the body of the post request : %s", err)
	}

	// create the new endpoint for the post request
	sid, err := getSID(body)
	if err != nil {
		return "", fmt.Errorf("error : %s", err)
	}

	return sid, nil
}

// return the result of a job get by its SID
func RetrieveJobResult(client *splunk.SplunkClient, spRequest *splunk.SplunkRequest) ([]map[string]string, error) {
	// make the get request
	getResp, err := GetJob(client, spRequest)
	if err != nil {
		return nil, fmt.Errorf("error while making the get request : %s", err)
	}
	// get the body of the response
	getBody, err := io.ReadAll(getResp.Body)

	if err != nil {
		return nil, fmt.Errorf("error while getting the body of the get request : %s", err)
	}

	// only get the result section of the response
	type Response struct {
		preview     bool
		init_offset int
		messages    []string
		fields      []map[string]string
		Results     []map[string]string `json:"results"`
	}

	results := Response{}
	errUmarshall := json.Unmarshal([]byte(getBody), &results)

	if errUmarshall != nil {
		return nil, errUmarshall
	}
	return results.Results, nil
}

// Return the sid from the body of the given response
func getSID(resp []byte) (string, error) {
	respJson := string(resp)

	var sid map[string]string
	errUmarshall := json.Unmarshal([]byte(respJson), &sid)

	if errUmarshall != nil {
		return "", errUmarshall
	}

	if len(sid) <= 0 {
		return "", fmt.Errorf("no sid found")
	}
	return sid["sid"], nil
}
