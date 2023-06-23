package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	logger "github.com/sirupsen/logrus"
)

// create an authentication key depending on the method provided
//
//	three methods are availables
//		1.  HTTP Authorization tokens : a seesion key
//		2.  Splunk Authentication tokens
//		3.  Basic Authentication
func CreateAuthenticationKey(client *SplunkClient) (string, error) {

	if client.Token != "" {
		if strings.HasPrefix(client.Token, "Splunk") {
			return "", fmt.Errorf("wrong authentication method. HTTP Authorization tokens used instead of Splunk Authentication tokens")
		}
		if strings.HasPrefix(client.Token, "Basic") {
			return "", fmt.Errorf("wrong authentication method. Basic Authentication used instead of Splunk Authentication tokens")
		}
		if !strings.HasPrefix(client.Token, "Bearer") {
			return "Bearer " + client.Token, nil
		}
	} else if client.SessionKey != "" {
		if strings.HasPrefix(client.SessionKey, "Bearer") {
			return "", fmt.Errorf("wrong authentication method. Splunk Authentication tokens used instead of HTTP Authorization tokens")
		}
		if strings.HasPrefix(client.SessionKey, "Basic") {
			return "", fmt.Errorf("wrong authentication method. Basic Authentication used instead of HTTP Authorization tokens")
		}
		if !strings.HasPrefix(client.SessionKey, "Splunk") {
			return "Splunk " + client.SessionKey, nil
		}
	} else if client.Username != "" && client.Password != "" {
		return base64.StdEncoding.EncodeToString([]byte(client.Username + ":" + client.Password)), nil
	}

	return "", fmt.Errorf("no authentication method provided")
}

// return the message of the error when got http error
func HandleHttpError(body []byte) (string, error) {

	var bodyJson map[string][]map[string]string
	errUmarshall := json.Unmarshal([]byte(body), &bodyJson)
	if errUmarshall != nil {
		return "", errUmarshall
	}

	if len(bodyJson) > 1 {
		return bodyJson["messages"][0]["text"], nil
	}
	return "", fmt.Errorf("incorrect format")
}

func MakeHttpRequest(client *SplunkClient, method string, spRequest *SplunkRequest, params url.Values) (*http.Response, error) {

	// create a new request
	req, err := http.NewRequest(method, client.Endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	// add the headers
	if spRequest.Headers == nil {
		spRequest.Headers = map[string]string{}
	}

	token, err := CreateAuthenticationKey(client)
	if err != nil {
		return nil, err
	}
	spRequest.Headers["Authorization"] = token
	logger.Infof("Authorization token: %s", token)
	logger.Infof("Authorization : %s", spRequest.Headers["Authorization"])
	for header, val := range spRequest.Headers {
		req.Header.Add(header, val)
	}

	// get the response
	resp, err := client.Client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func configureLogger(eventID, keptnContext string) {
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel, err := logger.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logger.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			logger.SetLevel(logLevel)
		}
	}
}
