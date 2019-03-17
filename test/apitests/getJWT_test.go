package main

import (
    "strings"
    "net/http"
    "io/ioutil"
    "testing"
)

func TestGetJWT(t *testing.T) {

    url := "http://localhost:3000/devices/login/"

    payload := strings.NewReader("{\"devName\":\"testDevice\", \"uName\":\"a\", \"pwd\":\"testDevice\"}")

    req, _ := http.NewRequest("GET", url, payload)

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("cache-control", "no-cache")
    req.Header.Add("Postman-Token", "701c7cf5-d94b-4333-a0e9-bacb521c7f91")

    res, _ := http.DefaultClient.Do(req)

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    t.Log(res)
    t.Log(string(body))

}
