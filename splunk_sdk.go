package splunksdk_go

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
)

const PATH_JOBS_V2 = "services/search/v2/jobs/"
const PATH_SAVED_SEARCHES = "saved/searches/"

type SplunkCreds struct {
	Host     string
	Port     string
	Token    string
	Endpoint string
}

type SplunkRequest struct {
	Client  *http.Client
	Headers map[string]string
	Params  RequestParams
}

type RequestParams struct {
	SearchQuery string
	OutputMode  string `default:"json"`
	ExecMode    string `default:"blocking"` // splunk returns a job SID only if the job is complete
}

// Return a metric from a new created job
func GetMetricFromNewJob(spRequest *SplunkRequest, spCreds *SplunkCreds) (float64, error) {
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel, err := logger.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logger.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			logger.SetLevel(logLevel)
		}
	}
	const RESULTS_URI = "results"

	endpoint, err := CreateJobEndpoint(spCreds)
	if err != nil {
		return -1, fmt.Errorf("error while creating the endpoint : %s", err)
	}

	spCreds.Endpoint = endpoint
	spRequest.Params.SearchQuery = validateSearchQuery(spRequest.Params.SearchQuery)
	sid, err := CreateJob(spRequest, spCreds)
	if err != nil {
		return -1, fmt.Errorf("error while creating the job : %s", err)
	}

	newEndpoint := endpoint + sid
	// check if the endpoint is correctly formed
	if !strings.HasSuffix(newEndpoint, "/") {
		newEndpoint += "/"
	}

	// the endpoint where to find the corresponding job
	spCreds.Endpoint = newEndpoint + RESULTS_URI

	res, err := RetrieveJobResult(spRequest, spCreds)
	if err != nil {
		return -1, fmt.Errorf("error while handling the results. Error message : %s", err)
	}
	// if the result is not a metric
	if len(res) != 1 {
		return -1, fmt.Errorf("incorrect search result. Error message : %v", err)
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

func post(spRequest *SplunkRequest, spCreds *SplunkCreds) (*http.Response, error) {

	return httpRequest("POST", spRequest, spCreds)
}

func get(spRequest *SplunkRequest, spCreds *SplunkCreds) (*http.Response, error) {

	return httpRequest("GET", spRequest, spCreds)
}

func validateSearchQuery(searchQuery string) string {
	// the search must start with the "search" keyword
	const QUERY_PREFIX = "search "
	if !strings.HasPrefix(searchQuery, QUERY_PREFIX) {
		return QUERY_PREFIX + searchQuery
	}
	return searchQuery
}

// this function create a new job and return its SID
func CreateJob(spRequest *SplunkRequest, spCreds *SplunkCreds) (string, error) {

	// params to send
	params := url.Values{}
	params.Add("search", spRequest.Params.SearchQuery)
	if spRequest.Params.OutputMode == "" {
		spRequest.Params.OutputMode = "json"
	}
	if spRequest.Params.ExecMode == "" {
		spRequest.Params.ExecMode = "blocking"
	}

	params.Add("output_mode", spRequest.Params.OutputMode)
	params.Add("exec_mode", spRequest.Params.ExecMode)

	logger.Infof("IN CREATE JOB PAR: %v", params)
	logger.Infof("IN CREATE JOB : %v", spRequest)
	resp, err := post(spRequest, spCreds)
	if err != nil {
		return "", fmt.Errorf("error while making the post request : %s", err)
	}
	body, err := io.ReadAll(resp.Body)
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
func RetrieveJobResult(spRequest *SplunkRequest, spCreds *SplunkCreds) ([]map[string]string, error) {
	// make the get request
	getResp, err := get(spRequest, spCreds)
	if err != nil {
		return nil, fmt.Errorf("error while making the get request : %s", err)
	}
	// get the body of the response
	getBody, err := io.ReadAll(getResp.Body)

	if err != nil {
		return nil, fmt.Errorf("error while getting the body of the get request : %s", err)
	}

	// only get the result section of the response
	var results map[string][]map[string]string
	json.Unmarshal([]byte(getBody), &results)

	return results["results"], nil
}

func CreateJobEndpoint(sc *SplunkCreds) (string, error) {
	host := sc.Host
	port := sc.Port

	match := `^((localhost)|(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})|([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}))$`

	if !regexp.MustCompile(match).MatchString(host) {
		return "", fmt.Errorf("")
	}

	return "https://" + net.JoinHostPort(host, port) + "/" + PATH_JOBS_V2, nil
}

// Return the sid from the body of the given response
func getSID(resp []byte) (string, error) {
	respJson := string(resp)

	var sid map[string]string
	json.Unmarshal([]byte(respJson), &sid)
	logger.Infof("IN GET SID : %s", resp)
	logger.Infof("IN GET SID2 : %s", sid)
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

func httpRequest(method string, spRequest *SplunkRequest, spCreds *SplunkCreds) (*http.Response, error) {

	// parameters of the request
	params := url.Values{}
	params.Add("search", spRequest.Params.SearchQuery)
	params.Add("output_mode", spRequest.Params.OutputMode)
	params.Add("exec_mode", spRequest.Params.ExecMode)

	// create a new request
	req, err := http.NewRequest(method, spCreds.Endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	// add the headers

	if spRequest.Headers == nil {
		spRequest.Headers = map[string]string{
			"Authorization": getBearer(spCreds.Token),
		}
	} else {
		for header, val := range spRequest.Headers {
			req.Header.Add(header, val)
		}

		req.Header.Set("Authorization", getBearer(spCreds.Token))
	}
	// get the response
	resp, err := spRequest.Client.Do(req)

	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()
	return resp, nil
}
