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
    "strings"
)


type DevEnv struct {
    db models.DeviceStore
}

var devEnv *DevEnv

type DevReq struct {
    NumEntries int `json:"numEntries"`
    Offset int `json:"offset"`
    ReqType string `json:"reqType"` // Can be "history", "heartbeat", "downlinkMsg"
    UplinkMsg string `json:"uplinkMsg"`
    PingTime int `json:"pingTime"`
    DownlinkMsg string `json:"downlinkMsg"`
    TunnelStatus string `json:"tunnelStatus"` // Can be "activate", "requested", "active", "closed"
}


// TODO : Make message structures more clear
type SSHReq struct {
    Port string `json:"port"`
    ReqType string `json:"reqType"`
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

func removeSSHKey(sshKey string) {

}


func DeviceManagerHandler(w http.ResponseWriter, r *http.Request) {
    emailCookie, err := r.Cookie("email")
    tableNameCookie, err1 := r.Cookie("dev_table")
    if err1 != nil && err != nil {
        if err == http.ErrNoCookie {
            http.Redirect(w, r, "/login/", http.StatusFound)
            return
        }
        http.Redirect(w, r, "/login/", http.StatusFound)
        return
    }
    email := emailCookie.Value
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
        email_tbl := strings.Replace(email,"@","_",-1)
        email_tbl = strings.Replace(email_tbl,".","_",-1)
        err = devEnv.db.CreateDeviceTable(email_tbl+"_"+devName) // TODO: Have a better way of creating device log
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
    usr, err := usrEnv.db.AUser(email)
    email_tbl := strings.Replace(usr.Email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    var deviceLogs []*models.DeviceLog
    if err != nil {
        log.Println(err)
    }
    if r.Method == "GET" {
        devName := r.URL.Query()["devName"][0]
        device, err := devEnv.db.ADevice(devTable, devName)
        if err != nil {
            log.Println(err)
        }
        deviceLogs, err = devEnv.db.GetDeviceLogs(email_tbl+"_"+device.DevName, 0, 10)
        if err != nil {
            log.Println(err)
        }
        //Ignore error
        if template, err := views.RenderDevicePreview(w, deviceLogs, device); err != nil {
            w.Write([]byte("No Data available"))
            return
        } else {
            w.Write(template)
            return
        }
        w.Write([]byte("Failed to get"))
    }
    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println(err)
        }
        devName := r.Form["devName"][0]
        downlinkMsg := r.Form["downlinkMsg"][0]
        reqType := r.Form["reqType"][0]
        tunnelStatus := r.Form["tunnelStatus"][0]
        downlinkPayload := &DevReq{ReqType:reqType,TunnelStatus:tunnelStatus, DownlinkMsg:downlinkMsg}

        email_tbl := strings.Replace(usr.Email,"@","_",-1)
        email_tbl = strings.Replace(email_tbl,".","_",-1)
        cacheId := email_tbl+"_"+devName
        err := devEnv.db.InsertDeviceDownlinkLog(email_tbl+"_"+devName, downlinkMsg, tunnelStatus)
        if err != nil {
            log.Println(err)
        }
        if reqType == "downlinkMsg" { // TODO : Implement queue here
            devDownlinkCache.Set(cacheId, downlinkPayload, cache.DefaultExpiration)
            w.WriteHeader(http.StatusOK)
        }
        if reqType == "ssh" {
            if tunnelStatus == "schedule" {
                device, err := devEnv.db.ADevice(devTable, devName)
                if err != nil {
                    log.Println("SSH Key not found")
                }
                port := AddKey(email, device.SSHKey)
                sshReq := &SSHReq{ReqType:reqType,Port:port}
                devDownlinkCache.Set(cacheId, sshReq, cache.DefaultExpiration)
                log.Println("Adding Key")
            }
            if tunnelStatus == "launch" {
            }
            if tunnelStatus == "stop" {
            }

        }
    }
}


func DeviceDataHandler(w http.ResponseWriter, r *http.Request, devClaims *DevClaims) {
    var devReq DevReq
    json.NewDecoder(r.Body).Decode(&devReq)
    email_tbl := strings.Replace(devClaims.Email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    if r.Method == "GET" {
    }

    if r.Method == "POST" {
        if devReq.ReqType == "heartbeat" {
            err := devEnv.db.InsertDeviceUplinkLog(email_tbl+"_"+devClaims.DevName, devReq.UplinkMsg, devReq.PingTime)
            if err != nil {
                log.Println(err)
                w.WriteHeader(http.StatusNotFound)
                return
            }
            cacheId := email_tbl+"_"+devClaims.DevName
            if val, found := devDownlinkCache.Get(cacheId); found {
                json.NewEncoder(w).Encode(val) 
                devDownlinkCache.Delete(cacheId)
                return
            }
            w.WriteHeader(http.StatusOK)
        }
        if devReq.ReqType == "uplinkMsg" {
            err := devEnv.db.InsertDeviceUplinkLog(devClaims.Email+"_"+devClaims.DevName, devReq.UplinkMsg, devReq.PingTime)
            if err != nil {
                log.Println(err)
                w.WriteHeader(http.StatusNotFound)
                return
            }
            cacheId := email_tbl+"_"+devClaims.DevName
            if val, found := devDownlinkCache.Get(cacheId); found {
                json.NewEncoder(w).Encode(val) 
                devDownlinkCache.Delete(cacheId)
                return
            }
            w.WriteHeader(http.StatusOK)
        }
        if devReq.ReqType == "downlinkMsg" { // TODO : Implement queue here
            downlinkPayload := &DevReq{ReqType:devReq.ReqType,
                                        TunnelStatus:devReq.TunnelStatus, DownlinkMsg:devReq.DownlinkMsg}
            cacheId := email_tbl+"_"+devClaims.DevName
            err := devEnv.db.InsertDeviceDownlinkLog(email_tbl+"_"+devClaims.DevName, devReq.DownlinkMsg, devReq.TunnelStatus)
            if err != nil {
                log.Println(err)
            }
            devDownlinkCache.Set(cacheId, downlinkPayload, cache.DefaultExpiration)
            w.WriteHeader(http.StatusOK)
        }

    }
}
