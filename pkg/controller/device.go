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


type UsrListDevsReq struct {
    NumEntries int `json:"numEntries"`
    Offset int `json:"offset"`
}

type DownlinkReq struct {
    DownlinkMsg string `json:"downlinkMsg"`
}

type DownlinkResp struct {
    Port string `json:"port"`
    DownlinkMsg string `json:"downlinkMsg"`
    RespType string `json:"reqType"` // schedule, stop or downlink
}

// Used my device and web, web ignores port
type SSHReq struct {
    TunnelStatus string `json:"tunnelStatus"` // schedule, launch, stop 
    Port string `json:"port"`
}

type UplinkReq struct {
    UplinkMsg string `json:"uplinkMsg"`
}

type HeartBeatReq struct {
    PingTime int `json:"pingTime"`
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


func getEmailTableName(email string) (string) {
    email_tbl := strings.Replace(email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    return email_tbl
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
    email_tbl := getEmailTableName(email)
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
        err = devEnv.db.CreateDeviceTable(email_tbl+"_"+devName) // TODO: Have a better way of creating device log
        if err != nil {
            log.Println(err)
        }
    }
    // TODO: Delete SSH key as well
    if r.Method == "DELETE" {
        devName := r.URL.Query().Get("devName")
        confDevName := r.URL.Query().Get("confDevName")
        if devName == confDevName {
            err = devEnv.db.DeleteDevice(email_tbl+"_"+tableName, devName)
            if err != nil {
                log.Println(err)
            }
            err = devEnv.db.DropDeviceTable(email_tbl+"_"+devName)
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


func GetDeviceInfo(w http.ResponseWriter, r *http.Request, email string, devTable string) {
    email_tbl := getEmailTableName(email)
    var deviceLogs []*models.DeviceLog
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
}

func GetDeviceSSHStatus(w http.ResponseWriter, r *http.Request, email string, devTable string) {
    email_tbl := getEmailTableName(email)
    if r.Method == "GET" {
        devName := r.URL.Query()["devName"][0]
        device, err := devEnv.db.ADevice(devTable, devName)
        if err != nil {
            log.Println("GetDeviceSSHStatus",err)
        }
        deviceStopLog, err := devEnv.db.GetLatestSSHStatus(email_tbl+"_"+device.DevName, "stop")
        if err != nil {
            log.Println("deviceStopLog",err)
        }
        deviceLaunchLog, err := devEnv.db.GetLatestSSHStatus(email_tbl+"_"+device.DevName, "launch")
        if err != nil {
            log.Println("deviceLaunchLog",err)
        }
        deviceScheduleLog, _ := devEnv.db.GetLatestSSHStatus(email_tbl+"_"+device.DevName, "schedule")
        if err != nil {
            log.Println("deviceScheduleLog",err)
        }
		if (deviceScheduleLog != nil) && (deviceLaunchLog != nil ) {
			if(deviceStopLog.LastSeen.After(deviceLaunchLog.LastSeen)) {
				val := &SSHReq{Port:"", TunnelStatus:"stopped"}
				json.NewEncoder(w).Encode(val)
				return
			}
			if(deviceLaunchLog.LastSeen.After(deviceScheduleLog.LastSeen)) {
				val := &SSHReq{Port:deviceLaunchLog.Port.String, TunnelStatus:"launch"}
				json.NewEncoder(w).Encode(val)
				return
			} else {
				val := &SSHReq{Port:"", TunnelStatus:"scheduled"}
				json.NewEncoder(w).Encode(val)
				return
			}
		} else {
				val := &SSHReq{Port:"", TunnelStatus:"sleep"}
				json.NewEncoder(w).Encode(val)
				return
		}

    }
}


func SendDeviceDownlink(w http.ResponseWriter, r *http.Request, email string, devTable string) {
    usr, err := usrEnv.db.AUser(email)
    email_tbl := strings.Replace(usr.Email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    if err != nil {
        log.Println(err)
    }
    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println(err)
        }
        devName := r.Form["devName"][0]
        downlinkMsg := r.Form["downlinkMsg"][0]
        email_tbl := getEmailTableName(usr.Email)
        cacheId := email_tbl+"_"+devName
        err := devEnv.db.InsertDeviceDownlinkLog(email_tbl+"_"+devName, downlinkMsg, "")
        if err != nil {
            log.Println(err)
        }
        downlinkPayload := &DownlinkResp{RespType:"downlink", DownlinkMsg:downlinkMsg}
        devDownlinkCache.Set(cacheId, downlinkPayload, cache.DefaultExpiration)
        w.WriteHeader(http.StatusOK)
    }
}

func SendSSHRequest(w http.ResponseWriter, r *http.Request, email string, devTable string) {
    email_tbl := getEmailTableName(email)
    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println(err)
        }
        devName := r.Form["devName"][0]
        downlinkMsg := r.Form["downlinkMsg"][0]
        tunnelStatus := r.Form["tunnelStatus"][0]
        cacheId := email_tbl+"_"+devName
        err := devEnv.db.InsertDeviceDownlinkLog(email_tbl+"_"+devName, downlinkMsg, tunnelStatus)
        if err != nil {
            log.Println(err)
        }
        device, err := devEnv.db.ADevice(devTable, devName)
        if tunnelStatus == "schedule" {
            if err != nil {
                log.Println("SSH Key not found")
            }
            err = DelDeviceKey(email_tbl, device.SSHKey)
            port := AddDeviceKey(email_tbl, device.SSHKey)
            sshReq := &DownlinkResp{RespType:tunnelStatus,Port:port}
            devDownlinkCache.Set(cacheId, sshReq, cache.DefaultExpiration)
        }
        if tunnelStatus == "launch" {
        }
        if tunnelStatus == "stop" {
            log.Println("Stopping SSH Tunnel")
            err = DelDeviceKey(email_tbl, device.SSHKey)
            log.Println("Stopped")
        }
    }
}

func SendHeartBeat(w http.ResponseWriter, r *http.Request, devClaims *DevClaims) {
    var heartBeatReq HeartBeatReq
    json.NewDecoder(r.Body).Decode(&heartBeatReq)
    email_tbl := getEmailTableName(devClaims.Email)
    if r.Method == "POST" {
        err := devEnv.db.InsertDeviceUplinkLog(email_tbl+"_"+devClaims.DevName, "", heartBeatReq.PingTime)
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
}

func SendUplink(w http.ResponseWriter, r *http.Request, devClaims *DevClaims) {
    var  uplinkReq UplinkReq
    json.NewDecoder(r.Body).Decode(&uplinkReq)
    email_tbl := getEmailTableName(devClaims.Email)
    if r.Method == "POST" {
        err := devEnv.db.InsertDeviceUplinkLog(email_tbl+"_"+devClaims.DevName, uplinkReq.UplinkMsg, -1)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusNotFound)
            return
        }
        cacheId := email_tbl+"_"+devClaims.DevName
        // Return back msg of typle UsrDownlinkResp
        if val, found := devDownlinkCache.Get(cacheId); found {
            json.NewEncoder(w).Encode(val)
            devDownlinkCache.Delete(cacheId)
            return
        }
        w.WriteHeader(http.StatusOK)
    }

}

func MakeTunnelRequest(w http.ResponseWriter, r *http.Request, devClaims *DevClaims) {
    var  sshReq SSHReq
    json.NewDecoder(r.Body).Decode(&sshReq)
    email_tbl := getEmailTableName(devClaims.Email)
    log.Println("Reached make tunnel request")
    if r.Method == "POST" {
        log.Println("Inserting Tunnel request log \t", sshReq)
        err := devEnv.db.InsertDeviceSSHLog(email_tbl+"_"+devClaims.DevName, sshReq.TunnelStatus, sshReq.Port)
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusNotFound)
            return
        }
        //TODO : Make port availability check, etc
        w.WriteHeader(http.StatusOK)
    }
}
