package main

import (
	"testing"
	"strings"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://localhost:3000/devices/data/"

	payload := strings.NewReader("{\"reqType\":\"heartbeat\", \"pingTime\":10}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("authToken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZOYW1lIjoidGVzdERldmljZTIiLCJwd2QiOiJ0ZXN0RGV2aWNlMiIsInVOYW1lIjoiYSJ9.WwpGe2ZYlCCYwi8CNxJHks7EHAZ3RK5PzzdJdZQwnsg")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "da53c175-0586-4b1f-8ba8-899306bda741")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	t.Log(res)
	t.Log(string(body))

}
