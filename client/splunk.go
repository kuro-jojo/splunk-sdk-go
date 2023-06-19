package client

import (
	"crypto/tls"
	"net/http"
)

type SplunkRequest struct {
	Headers map[string]string
	Params  RequestParams
}

type RequestParams struct {
	// splunk search in spl syntax
	SearchQuery string
	OutputMode  string `default:"json"`
	// splunk returns a job SID only if the job is complete
	ExecMode string `default:"blocking"`
	// earliest (inclusive) time bounds for the search
	EarliestTime string
	// latest (exclusive) time bounds for the search
	LatestTime string
}

type SplunkClient struct {
	Client     *http.Client
	Host       string
	Port       string
	Endpoint   string
	Token      string
	Username   string
	Password   string
	SessionKey string
	// if true, ssl verification is skipped
	SkipSSL bool
}

// create a new Client
func NewClient(client *http.Client, host string, port string, token string, username string, password string, sessionKey string, skipSSL bool) *SplunkClient {
	if skipSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return &SplunkClient{
		Client:     client,
		Host:       host,
		Port:       port,
		Token:      token,
		Username:   username,
		Password:   password,
		SessionKey: sessionKey,
		SkipSSL:    skipSSL,
	}
}

// create a new client that could connect with authentication tokens
func NewClientAuthenticatedByToken(client *http.Client, host string, port string, token string, skipSSL bool) *SplunkClient {
	if skipSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return &SplunkClient{
		Client:     client,
		Host:       host,
		Port:       port,
		Token:      token,
		Username:   "",
		Password:   "",
		SessionKey: "",
		SkipSSL:    skipSSL,
	}
}

// create a new client that could connect with authentication sessionKey
func NewClientAuthenticatedBySessionKey(client *http.Client, host string, port string, sessionKey string, skipSSL bool) *SplunkClient {
	if skipSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return &SplunkClient{
		Client:     client,
		Host:       host,
		Port:       port,
		SessionKey: sessionKey,
		Token:      "",
		Username:   "",
		Password:   "",
		SkipSSL:    skipSSL,
	}
}

// create a new client with basic authentication method
func NewBasicAuthenticatedClient(client *http.Client, host string, port string, username string, password string, skipSSL bool) *SplunkClient {
	if skipSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return &SplunkClient{
		Client:     client,
		Host:       host,
		Port:       port,
		Username:   username,
		Password:   password,
		Token:      "",
		SessionKey: "",
		SkipSSL:    skipSSL,
	}
}
