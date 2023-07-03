package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(client.Username+":"+client.Password)), nil
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
		spRequest.Headers= map[string]string{}
	}

	token, err := CreateAuthenticationKey(client)
	if err != nil {
		return nil, err
	}
	spRequest.Headers= map[string]string{"Authorization":token}

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

func MakeAlertHttpRequest(client *SplunkClient, method string, spRequest *SplunkAlert, params url.Values) (*http.Response, error) {

	// create a new request
	req, err := http.NewRequest(method, client.Endpoint, strings.NewReader(params.Encode()))
	
	if err != nil {
		return nil, err
	}
	// add the headers
	if spRequest.Headers == nil {
		spRequest.Headers= map[string]string{}
	}

	token, err := CreateAuthenticationKey(client)
	if err != nil {
		return nil, err
	}
	spRequest.Headers= map[string]string{"Authorization":token}

	for header, val := range spRequest.Headers {
		req.Header.Add(header, val)
	}
	fmt.Printf("--> %s\n\n", formatRequest(req))
	// get the response
	resp, err := client.Client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
	  name = strings.ToLower(name)
	  for _, h := range headers {
		request = append(request, fmt.Sprintf("%v: %v", name, h))
	  }
	}
	
	// If this is a POST, add post data
	if r.Method == "POST" {
	   r.ParseForm()
	   request = append(request, "\n")
	   request = append(request, r.Form.Encode())
	} 
	 // Return the request as a string
	 return strings.Join(request, "\n")
   }