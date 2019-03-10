package main

import (
    "testing"
	"strings"
	"net/http"
	"io/ioutil"
)

func TestSendDownlink(t *testing.T) {

	url := "http://localhost:3000/devices/data/"

	payload := strings.NewReader("{\"reqType\":\"downlinkMsg\", \"activateSSH\":true, \"downlinkMsg\":\"Test message e\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("authToken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZOYW1lIjoidGVzdERldmljZTIiLCJwd2QiOiJ0ZXN0RGV2aWNlMiIsInVOYW1lIjoiYSJ9.WwpGe2ZYlCCYwi8CNxJHks7EHAZ3RK5PzzdJdZQwnsg")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "a8e66299-b7c1-482c-94ac-6469d650fdcb")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	t.Log(res)
	t.Log(string(body))

}
