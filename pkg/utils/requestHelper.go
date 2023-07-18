package utils

import (
	"net"
	"strings"

	splunk "github.com/kuro-jojo/splunk-sdk-go/src/client"
)

func ValidateSearchQuery(searchQuery string) string {
	// the search must start with the "search" keyword
	const query_prefix = "search "
	if !strings.HasPrefix(searchQuery, query_prefix) {
		return query_prefix + searchQuery
	}
	return searchQuery
}
func CreateEndpoint(client *splunk.SplunkClient, service string) {
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
