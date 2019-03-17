package main

import (
	"testing"
	"strings"
	"net/http"
	"io/ioutil"
)

func TestSendUplink(t *testing.T) {

	url := "http://localhost:3000/devices/data/"

	payload := strings.NewReader("{\"reqType\":\"heartbeat\", \"pingTime\":10}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("authToken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZOYW1lIjoidGVzdERldmljZSIsInB3ZCI6InRlc3REZXZpY2UiLCJ1TmFtZSI6ImEifQ.g-Fi7EVwty1QAowDabtgJPLtzvyqt2O5XHaXP5QYUGM")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "da53c175-0586-4b1f-8ba8-899306bda741")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	t.Log(res)
	t.Log(string(body))

}
