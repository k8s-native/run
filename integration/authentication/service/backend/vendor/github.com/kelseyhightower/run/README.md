# run

[![GoDoc](https://godoc.org/github.com/kelseyhightower/run?status.svg)](https://pkg.go.dev/github.com/kelseyhightower/run)

The run package provides a set of Cloud Run helper functions and does not leverage any third party dependencies.

## Usage

```Go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/kelseyhightower/run"
)

func main() {
    // Set up a logger that generates structured logs optimized for
    // Cloud Run.
    logger, err := run.NewLogger()
    if err != nil {
        log.Fatal(err)
    }

    logger.Notice("Starting helloworld service...")

    // Easy access to secrets stored in Secret Manager.
    secret, err := run.AccessSecret("foo")
    if err != nil {
        logger.Error(err)
        os.Exit(1)
    }

    _ = secret

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello world!\n"))
    })

    // Start an HTTP server listening on the address defined
    // by the Cloud Run container runtime contract.
    log.Fatal(run.ListenAndServe(nil))
}
```

### Service Authentication

run takes the pain out of [service-to-service authentication](https://cloud.google.com/run/docs/authenticating/service-to-service)

```Go
package main

import (
    "log"
    "net/http"

    "github.com/kelseyhightower/run"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        request, err := http.NewRequest("GET", "https://example-6bn2iswfgq-uw.a.run.app", nil)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        // Use the run.Transport to automatically attach ID tokens to outbound requests
        // and optionally expand service names using the Cloud Run API.
        // See https://pkg.go.dev/github.com/kelseyhightower/run?tab=doc#Transport
        client := http.Client{Transport: &run.Transport{EnableServiceNameResolution: false}}

        response, err := client.Do(request)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        defer response.Body.Close()

        // Do something with response.
    })

    log.Fatal(run.ListenAndServe(nil))
}
```

## Status

This package is experimental and should not be used or assumed to be stable. Breaking changes are guaranteed to happen.
