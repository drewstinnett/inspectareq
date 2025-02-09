/*
Package main is the executable for the enable-env example
*/
package main

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/drewstinnett/inspectareq"
)

func main() {
	fmt.Printf(`Use %v to enable curl printing, and %v to enable httpie printing`, inspectareq.CurlEnv, inspectareq.HTTPieEnv)
	reqBody := bytes.NewBufferString(`{"username": "alice", "password": "secret"}`)
	req, err := http.NewRequest("POST", "https://pie.dev/anything", reqBody)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add some headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Debug", "true")
	req.Header.Set("Authorization", "Bearer my-token")

	if err := processRequest(req); err != nil {
		log.Fatalf("error processing request: %v", err)
	}
}

func processRequest(req *http.Request) error {
	slog.Info("debbuging status", "enabled", inspectareq.Enabled())

	if err := inspectareq.Print(req); err != nil {
		log.Fatalf("error printing request: %v", err)
	}

	return nil
}
