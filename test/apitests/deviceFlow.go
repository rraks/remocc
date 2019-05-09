package main

import (
	"strings"
	"net/http"
    "time"
    "log"
    "io"
    "encoding/json"
    "fmt"
    "bytes"
)

var uplinkUrl = "http://localhost:3000/devices/data/uplink/"
var sshUrl = "http://localhost:3000/devices/data/ssh/"

type Jwtresp struct {
    Token string `json:"token"`
}

func GetJWT(devName string, email string, pwd string) (string){
    var jwtresp Jwtresp
    url := "http://localhost:3000/devices/login/"
    payload := strings.NewReader("{\"devName\":\""+devName+"\", \"email\":\""+email+"\", \"pwd\":\""+pwd+"\"}")
    req, _ := http.NewRequest("GET", url, payload)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("cache-control", "no-cache")
    req.Header.Add("Postman-Token", "701c7cf5-d94b-4333-a0e9-bacb521c7f91")
    res, _ := http.DefaultClient.Do(req)
    log.Println(res)
    json.NewDecoder(res.Body).Decode(&jwtresp)
    return jwtresp.Token
}

func SendUplink(payload io.Reader, jwt string) (map[string]interface{}) {
    var response map[string]interface{}
	req, _ := http.NewRequest("POST", uplinkUrl, payload)
	req.Header.Add("authToken", jwt)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)
    json.NewDecoder(res.Body).Decode(&response)
    return response
}

func SendSSHReq(port string, jwt string) (map[string]interface{}) {
    log.Println("SendSSHReq")
    var response map[string]interface{}
    payload := make(map[string]string)
    payload["port"] = port
    payload["tunnelStatus"] = "launch"
    payloadStr, err := json.Marshal(payload)
    if err != nil {
        log.Println(err)
    }
    log.Println("Sending paylod \t", string(payloadStr))
	req, _ := http.NewRequest("POST", sshUrl, bytes.NewReader(payloadStr))
	req.Header.Add("authToken", jwt)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)
    json.NewDecoder(res.Body).Decode(&response)
    log.Println("Send Response")
    return response
}


func main() {
    var devName = "testDevice"
    var email = "a@a.com"
    var devPassword = "a@a.com"
    payload := strings.NewReader("{\"uplinkMsg\":\"yada\", \"pingTime\":10}")
    jwt := GetJWT(devName, email, devPassword)
    for {
        resp := SendUplink(payload, jwt)
        log.Println("[UPLINK] \t Sent \t", payload)
        log.Println("[UPLINK] \t Received \t ",resp)
        if resp != nil {
            reqType := resp["reqType"]
            if(reqType == "schedule") {
                log.Println("[DOWNLINK] \t SSH \t", resp)
                port := fmt.Sprintf("%v", resp["port"])
                log.Println("Port is ", port)
                resp = SendSSHReq(port, jwt)
                if resp != nil {
                    log.Println(resp)
                }
                log.Println("[END]")
            }
        }
        time.Sleep(5*time.Second)
    }
}
