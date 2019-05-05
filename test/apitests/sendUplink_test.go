package main

import (
	"testing"
	"strings"
	"net/http"
	"io/ioutil"
)

func TestSendUplink(t *testing.T) {

	url := "http://localhost:3000/devices/data/heartbeat/"

	payload := strings.NewReader("{\"reqType\":\"heartbeat\", \"pingTime\":10}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("authToken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZOYW1lIjoidGVzdERldmljZSIsImVtYWlsIjoiYUBhLmNvbSIsInB3ZCI6InRlc3REZXZpY2UifQ.ugYStYg2HBCyg2DyvGgbmZrofBdVWLlY5HQcX_Q12gI")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	t.Log(res)
	t.Log(string(body))

}
