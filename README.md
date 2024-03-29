# Splunk SDK for Go

_Version 1.7.0_

The Splunk Software Development Kit for Go contains functions designed to enable developers to communicate with Splunk Enterprise through the splunk API.

:warning: This version is more focused on getting metrics from Splunk Enterprise.

## Getting started with the Splunk SDK for Go

### Requirements

Here's what you need to get going with the Splunk Enterprise SDK for Go.

- Go 1.18+

  The Splunk SDK for Go has been tested with Go version 1.18 to 1.20

- Spunk Enterprise 9.0.4

      The Splunk SDK has been tested with Splunk Enterprise 9.0.4

  If you haven't already installed Splunk Enterprise, download it [here](http://www.splunk.com/download). For more information, see the Splunk Enterprise Installation Manual.

### Install the SDK

Use the following command to install the Splunk SDK for Go

    go get -u github.com/kuro-jojo/splunk-sdk-go

### Example [![Go Reference](https://pkg.go.dev/badge/github.com/kuro-jojo/splunk-sdk-go.svg)](https://pkg.go.dev/github.com/kuro-jojo/splunk-sdk-go)

You'll need at first a Splunk enterprise instance running. If you don't have one, you can run a local instance with a docker image.

- For that you'll need _docker_ to be installed
- Then run a local Splunk enterprise instance (check it on [docker](https://hub.docker.com/r/splunk/splunk)) :

         docker run -d -p 8000:8000 -e "SPLUNK_START_ARGS=--accept-license" -e "SPLUNK_PASSWORD=<password>" --name splunk splunk/splunk:latest

  After the container starts up successfully and enters the "healthy" state, you should be able to access SplunkWeb at http://localhost:8000 with admin:\<password>.

#### Following are the different ways to connect to Splunk Enterprise

> :warning: Avoid writing your sensitive information in your code in production. Use environment variables or a configuration file instead.

**Using username and password**

```go
    import (
        splunk "github.com/kuro-jojo/splunk-sdk-go/client"
    )
    ...
        splunkInstance := "localhost" // or your splunk instance IP
        splunkServerPort := "8089" // by default
        splunkUsername := "admin"
        splunkPassword := "myComplexPassword" //
        client := splunk.NewBasicAuthenticatedClient(
        &http.Client{
            Timeout: time.Duration(60) * time.Second,
        },
        splunkInstance,
        splunkServerPort,
        splunkUsername,
        splunkPassword,
        true, // if true : SSL verification disabled
```

**Using token authentication**

```go
    import (
        splunk "github.com/kuro-jojo/splunk-sdk-go/client"
    )
    ...
        splunkInstance := "localhost"
        splunkServerPort := "8089" // by default
        splunkToken := "eyJraWQiOiJzcGx1bmsuc2VjcmV0IiwiYWxnIjo..."
        client := splunk.NewClientAuthenticatedByToken(
        &http.Client{
            Timeout: time.Duration(60) * time.Second,
        },
        splunkInstance,
        splunkServerPort,
        splunkToken,
        true, // if true : SSL verification disabled
```

**Using authentication sessionKey**

```go
    import (
        splunk "github.com/kuro-jojo/splunk-sdk-go/client"
    )
    ...
        splunkInstance := "localhost"
        splunkServerPort := "8089" // by default
        splunkSessionKey := "ff8be3be-ef07-4576-..."
        client := splunk.NewClientAuthenticatedBySessionKey(
        &http.Client{
            Timeout: time.Duration(60) * time.Second,
        },
        splunkInstance,
        splunkServerPort,
        splunkSessionKey,
        true, // if true : SSL verification disabled
```

#### Create a new job

```go
...
import (
    splunk "github.com/kuro-jojo/splunk-sdk-go/client"
    job "github.com/kuro-jojo/splunk-sdk-go/jobs"
    ...
)
...
    // create the parameters for the search
    searchParameters := splunk.SearchParams{
        SearchQuery: "index=main | head 10",
    }

    spReq := splunk.SearchRequest{
        Params: searchParameters,
        Headers: map[string]string{
            "Content-Type": "application/text",
            "..."
        },
    }

    // create the job and get the sid of the job which will be used to get the results
    sid, err := job.CreateJob(client, &spReq)

    if err != nil {
        fmt.Printf("Got an error : %s", err)
        return
    }

```

#### Retrieving the results of a job

```go
...
import (
    splunk "github.com/kuro-jojo/splunk-sdk-go/client"
    job "github.com/kuro-jojo/splunk-sdk-go/jobs"
    ...
)
...

    // get the results of the search using the sid of the job
    results, err := job.RetrieveJobResult(client, sid)

    if err != nil {
        fmt.Printf("Got an error : %s", err)
        return
    }

```

#### Getting metric from a job

```go
...
import (
    splunk "github.com/kuro-jojo/splunk-sdk-go/client"
    job "github.com/kuro-jojo/splunk-sdk-go/jobs"
    ...
)
...
    // create the parameters for the search

    spReq := splunk.SearchRequest{
        Params: splunk.SearchParams{
            SearchQuery: "index=main | stats count",
        },
    }

    metric, err := job.GetMetricFromNewJob(client, &spReq)
    fmt.Println(metric)
    if err != nil {
        fmt.Printf("Got an error : %s", err)
        return
    }

```

## License

The Splunk Enterprise Software Development Kit for Go is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
