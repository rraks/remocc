package controller

//TODO :[
//          Verify uniqueness of name ]

import (
    "net/http"
    "github.com/rraks/remocc/pkg/models"
    "github.com/rraks/remocc/pkg/views"
    "log"
    "encoding/json"
    "github.com/patrickmn/go-cache"
    "time"
)


type DevEnv struct {
    db models.DeviceStore
}

var devEnv *DevEnv

type DevReq struct {
    NumEntries int `json: "numEntries"`
    Offset int `json: "offset"`
    ReqType string `json: "reqType"` // Can be "history", "heartbeat", "downlinkMsg"
    UplinkMsg string `json: "uplinkMsg"`
    PingTime int `json: "pingTime"`
    DownlinkMsg string `json: "downlinkMsg"`
    ActivateSSH bool `json: "activateSSH"`

}



var devDownlinkCache *cache.Cache


func init() {
    db, err := models.InitDB()
    if err != nil {
        panic(err)
    }
    devEnv = &DevEnv{db}

    //Create simple downlink cache
    devDownlinkCache = cache.New(2*time.Hour,4*time.Hour)
}

func DeviceManagerHandler(w http.ResponseWriter, r *http.Request) {
    tableNameCookie, err := r.Cookie("dev_table")
    if err != nil {
        if err == http.ErrNoCookie {
            http.Redirect(w, r, "/login/", http.StatusFound)
            return
        }
        http.Redirect(w, r, "/login/", http.StatusFound)
        return
    }
    tableName := tableNameCookie.Value
    err = r.ParseForm()
    if err != nil{
        http.Redirect(w, r, "/", http.StatusFound) // TODO: Flash error message here
    }

    if r.Method == "POST" {
        devName := r.Form["devName"][0]
        macId := r.Form["macId"][0]
        description := r.Form["devDescr"][0]
        sshKey := r.Form["sshKey"][0]
        devPwd := r.Form["devPwd"][0]
        devPwdHash,_ := HashPassword(devPwd)
        _, err = devEnv.db.NewDevice(tableName, devName, macId, description, sshKey, devPwdHash)
        if err != nil {
            log.Println(err)
        }
        err = devEnv.db.CreateDeviceTable(devName) // TODO: Have a better way of creating device log
        if err != nil {
            log.Println(err)
        }
    }
    if r.Method == "DELETE" {
        devName := r.URL.Query().Get("devName")
        confDevName := r.URL.Query().Get("confDevName")
        if devName == confDevName {
            err = devEnv.db.DeleteDevice(tableName, devName)
            if err != nil {
                log.Println(err)
            }
            err = devEnv.db.DropDeviceTable(devName)
            if err != nil {
                log.Println(err)
            }
        }
    }
    http.Redirect(w, r, "/", http.StatusFound) // TODO: Flash error message here
}


// TODO: Add a more qualified user context here instead of tableName
func FrontPageHandler(w http.ResponseWriter, r *http.Request, email string, tableName string) {
    allDevices, err := devEnv.db.AllDevices(tableName)
    if err != nil {
        log.Println(err)
    }
    views.RenderTableRow(w, allDevices) // Ignore error
}


func UserDataHandler(w http.ResponseWriter, r *http.Request, email string, devTable string) {
    if r.Method == "GET" {
        devName := r.URL.Query()["devName"][0]
        device, err := devEnv.db.ADevice(devTable, devName)
        if err != nil {
            log.Println(err)
        }
        deviceLogs, err := devEnv.db.GetDeviceLogs(device, 0, 10)
        if err != nil {
            log.Println(err)
        }
        if template, err := views.RenderDevicePreview(w, deviceLogs, device); err != nil {
            w.Write([]byte("No Data available"))
            return
        } else {
            w.Write(template)
            return
        }
        w.Write([]byte("Failed to get"))
    }
}


func DeviceDataHandler(w http.ResponseWriter, r *http.Request, devClaims *DevClaims) {
    var devReq DevReq
    json.NewDecoder(r.Body).Decode(&devReq)

    if r.Method == "GET" {
    }

    if r.Method == "POST" {
        if devReq.ReqType == "heartbeat" {
            err := devEnv.db.InsertDeviceLog(devClaims.DevName, devReq.UplinkMsg, devReq.PingTime)
            if err != nil {
                log.Println(err)
                w.WriteHeader(http.StatusNotFound)
                return
            }
            cacheId := devClaims.UName+"_"+devClaims.DevName
            if val, found := devDownlinkCache.Get(cacheId); found {
                json.NewEncoder(w).Encode(val) 
                devDownlinkCache.Delete(cacheId)
                return
            }
            w.WriteHeader(http.StatusOK)
        }
        if devReq.ReqType == "downlinkMsg" { // TODO : Implement queue here
            downlinkPayload := &DevReq{ReqType:devReq.ReqType,
                                        ActivateSSH:devReq.ActivateSSH, DownlinkMsg:devReq.DownlinkMsg}
            cacheId := devClaims.UName+"_"+devClaims.DevName
            devDownlinkCache.Set(cacheId, downlinkPayload, cache.DefaultExpiration)
            w.WriteHeader(http.StatusOK)
        }

    }
}

