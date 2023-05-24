package splunksdk

import (
	// "crypto/tls"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestSplunkAPI(t *testing.T) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	err := godotenv.Load("go.env")
	if os.Getenv("ENV") == "dev" {
		err = godotenv.Load(".env.local")
	}
	if err != nil {
		log.Fatalf("Error loading environment variables file")
	}
	apiToken := os.Getenv("SPLUNK_TOKEN")

	search := "source=/opt/splunk/var/log/secure.log sourcetype=osx_secure |stats count"
	// create the http client
	client := &http.Client{
		Transport: tr,
		Timeout: time.Duration(1) * time.Second,
	}
	sc := SplunkCreds{
		Host:  os.Getenv("SPLUNK_HOST"),
		Port:  os.Getenv("SPLUNK_PORT"),
		Token: apiToken,
		Endpoint: "",
	}

	// get the metric we want
	metric, err := GetMetric(client, &sc, search, nil)
	fmt.Printf("Endpoint : %s\n\n", sc)
	if err != nil {
		t.Errorf("Error : %s\n", err)
		return
	}
	var expectedMetric float64 = 9829
	if expectedMetric != metric {
		t.Errorf("Expected %v but got %v.", expectedMetric, metric)
	}
}
