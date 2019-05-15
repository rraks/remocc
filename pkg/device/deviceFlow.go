// Package device creates a binary for a device
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
    "os/exec"
)

var loginUrl = "http://webdev:3000/devices/login/"
var uplinkUrl = "http://webdev:3000/devices/data/uplink/"
var sshUrl = "http://webdev:3000/devices/data/ssh/"

var devName = "testDevice"
var email = "a@a.com"
var devPassword = "testDevice"

type Jwtresp struct {
    Token string `json:"token"`
}

func checkErr(err error, context string) {
    if err != nil {
        log.Println("[ERROR] \t", context)
        log.Println("\t\t", err)
    }
}

func exe_cmd(cmd string, args ...string) ([]byte) {
    log.Println("Sending command \n", cmd, args)
    out, err := exec.Command(cmd, args...).Output()
    if err != nil {
      log.Println("Exec Failed")
      log.Println(err)
    }
    return out
}

func GetJWT(devName string, email string, pwd string) (string){
    var jwtresp Jwtresp
    payload := strings.NewReader("{\"devName\":\""+devName+"\", \"email\":\""+email+"\", \"pwd\":\""+pwd+"\"}")
    req, _ := http.NewRequest("GET", loginUrl, payload)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("cache-control", "no-cache")
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
    if res, _ := http.DefaultClient.Do(req); res != nil {
        json.NewDecoder(res.Body).Decode(&response)
        return response
    } else {
        return nil
    }
}

func SendSSHReq(port string, jwt string) bool  {
    var response map[string]interface{}
    payload := make(map[string]string)
    payload["port"] = port
    payload["tunnelStatus"] = "launch"
    payloadStr, err := json.Marshal(payload)
    checkErr(err, "Unabled to marshall ssh request")
    req, _ := http.NewRequest("POST", sshUrl, bytes.NewReader(payloadStr))
    req.Header.Add("authToken", jwt)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("cache-control", "no-cache")
    res, _ := http.DefaultClient.Do(req)
    json.NewDecoder(res.Body).Decode(&response)
    if res.StatusCode == http.StatusOK {
        return true
    }
    return false
}


func runloop() {
    jwt := GetJWT(devName, email, devPassword)
    linux_usr := strings.Replace(email, "@", "_", -1)
    linux_usr = strings.Replace(linux_usr, ".", "_", -1)
    for {
        payload := strings.NewReader("{\"uplinkMsg\":\"yada\", \"pingTime\":10}")
        log.Println("[UPLINK] \t Sent \t", payload)
        if resp := SendUplink(payload, jwt); resp != nil {
            log.Println("[UPLINK] \t Received \t ",resp)
            reqType := resp["reqType"]
            if(reqType == "schedule") {
                log.Println("[DOWNLINK] \t SSH \t", resp)
                port := fmt.Sprintf("%v", resp["port"])
                log.Println("Port is ", port)
                if SendSSHReq(port, jwt){
                    go exe_cmd("nohup", "ssh", "-p2222", "-N", "-R", port+":localhost:"+"22", linux_usr+"@webdev", "> /dev/null 2>&1", "&")
                    log.Println("Starting SSH Session")
                }
                log.Println("[END]")
            }
        }
        time.Sleep(5*time.Second)
    }
}


func main() {
    runloop()
}
