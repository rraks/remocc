package controller

// device.go
// Provides structures and functions for device operations

import (
    "net/http"
    "github.com/rraks/remocc/pkg/models"
    "github.com/rraks/remocc/pkg/views"
    "log"
    "encoding/json"
    "github.com/patrickmn/go-cache"
    "time"
    "strings"
    "github.com/rraks/remocc/pkg/ssh"
    "github.com/rraks/remocc/pkg/scheduler"
)

// Downlink cache for all devices
var devDownlinkCache *cache.Cache

// Database environment for "device" related database accesses
type DevEnv struct {
    db models.DeviceStore
}
var devEnv *DevEnv


// Structure to list user's devices
type UsrListDevsReq struct {
    NumEntries int `json:"numEntries"`
    Offset int `json:"offset"`
}


// Structure users use to give either of the following - 
//  send a downlink message OR ssh schedule/stop request
type DownlinkReq struct {
    Port string `json:"port"`
    DownlinkMsg string `json:"downlinkMsg"`
    RespType string `json:"reqType"` // schedule, stop or downlink
}

// Structure users use to get SSH tunnel Status 
// Structure devices use to make an SSH launch request
type SSHReq struct {
    TunnelStatus string `json:"tunnelStatus"` // schedule, launch, stop 
    Port string `json:"port"`
}

// Structure devices use to send a message to server
type UplinkReq struct {
    UplinkMsg string `json:"uplinkMsg"`
}

// Structure devices use register that they are alive
type HeartBeatReq struct {
    PingTime int `json:"pingTime"`
}


// Initialize db and downlink cache
func init() {
    db, err := models.InitDB()
    if err != nil {
        panic(err)
    }
    devEnv = &DevEnv{db}
    //Create simple downlink cache
    devDownlinkCache = cache.New(2*time.Hour,4*time.Hour)
}


// Log errors
func checkErr(err error, context string) {
    if err != nil {
        log.Println("[ERROR] \t", context)
        log.Println("\t\t", err)
    }
}

// Make email name psql friendly
func getEmailTableName(email string) (string) {
    email_tbl := strings.Replace(email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    return email_tbl
}

// Function to manage devices. Handles device registration and deletion
func DeviceManagerHandler(w http.ResponseWriter, r *http.Request, email string, tableName string) {
    email_tbl := getEmailTableName(email)
    err := r.ParseForm()
    if err != nil {
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
        checkErr(err, "Passwords don't match")
        err = devEnv.db.CreateDeviceTable(email_tbl+"_"+devName) // TODO: Have a better way of creating device log
        checkErr(err, "Couldn't create device table")
    }
    // TODO: Delete SSH key as well
    if r.Method == "DELETE" {
        devName := r.URL.Query().Get("devName")
        confDevName := r.URL.Query().Get("confDevName")
        if devName == confDevName {
            err = devEnv.db.DeleteDevice("devices_" + email_tbl, devName)
            checkErr(err, "No such device exists")
            err = devEnv.db.DropDeviceTable(email_tbl+"_"+devName)
            checkErr(err, "No such device table exists")
        }
    }
    http.Redirect(w, r, "/", http.StatusFound) // TODO: Flash error message here
}


// Render template for front page displaying all the devices registered to a particular user
// TODO: Add a more qualified user context here instead of tableName
func FrontPageHandler(w http.ResponseWriter, r *http.Request, email string, tableName string) {
    allDevices, err := devEnv.db.AllDevices(tableName)
    checkErr(err, "Couldn't retreive device table")
    views.RenderTableRow(w, allDevices) // Ignore error
}


// Device preview page showing tunnel status and latest logs (up to 10)
func GetDeviceInfo(w http.ResponseWriter, r *http.Request, email string, devTable string) {
    email_tbl := getEmailTableName(email)
    var deviceLogs []*models.DeviceLog
    if r.Method == "GET" {
        devName := r.URL.Query()["devName"][0]
        device, err := devEnv.db.ADevice(devTable, devName)
        checkErr(err, "No such device exists")
        deviceLogs, err = devEnv.db.GetDeviceLogs(email_tbl+"_"+device.DevName, 0, 10)
        checkErr(err, "Couldn't retreive device logs")
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

// Get SSHStatus, i.e, if user has scheduled a tunnel to be made, if device has launched that tunnel, 
// or if it is stopped
func GetDeviceSSHStatus(w http.ResponseWriter, r *http.Request, email string, devTable string) {
    email_tbl := getEmailTableName(email)
    if r.Method == "GET" {
        devName := r.URL.Query()["devName"][0]
        device, err := devEnv.db.ADevice(devTable, devName)
        checkErr(err, "Couldn't retreuve device")
        deviceStopLog, err := devEnv.db.GetLatestSSHStatus(email_tbl+"_"+device.DevName, "stop")
        checkErr(err, "No stop logs")
        deviceLaunchLog, err := devEnv.db.GetLatestSSHStatus(email_tbl+"_"+device.DevName, "launch")
        checkErr(err, "No launch logs")
        deviceScheduleLog, _ := devEnv.db.GetLatestSSHStatus(email_tbl+"_"+device.DevName, "schedule")
        checkErr(err, "No schedule logs")
        _ = deviceStopLog // TODO: Do somethings with deviceStopLog
		if (deviceScheduleLog != nil) && (deviceLaunchLog == nil ) {
				val := &SSHReq{Port:"", TunnelStatus:"scheduled"}
				json.NewEncoder(w).Encode(val)
                return
        }
		if (deviceScheduleLog != nil) && (deviceLaunchLog != nil ) {
			if(deviceLaunchLog.LastSeen.After(deviceScheduleLog.LastSeen)) {
				val := &SSHReq{Port:deviceLaunchLog.Port.String, TunnelStatus:"launch"}
				json.NewEncoder(w).Encode(val)
				return
			} else {
				val := &SSHReq{Port:"", TunnelStatus:"scheduled"}
				json.NewEncoder(w).Encode(val)
				return
			}
		}
        if (deviceStopLog != nil) && (deviceLaunchLog != nil) {
			if(deviceStopLog.LastSeen.After(deviceLaunchLog.LastSeen)) {
				val := &SSHReq{Port:"", TunnelStatus:"stopped"}
				json.NewEncoder(w).Encode(val)
				return
            }
        }
    }
}


// Function users use to send a downlink message to the device
func SendDeviceDownlink(w http.ResponseWriter, r *http.Request, email string, devTable string) {
    usr, err := usrEnv.db.AUser(email)
    email_tbl := strings.Replace(usr.Email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    checkErr(err, "No such user exists")
    if r.Method == "POST" {
        if err := r.ParseForm(); err != nil {
            log.Println(err)
        }
        devName := r.Form["devName"][0]
        downlinkMsg := r.Form["downlinkMsg"][0]
        email_tbl := getEmailTableName(usr.Email)
        cacheId := email_tbl+"_"+devName
        err := devEnv.db.InsertDeviceDownlinkLog(email_tbl+"_"+devName, downlinkMsg, "")
        checkErr(err, "Couldn't insert downlink log")
        downlinkPayload := &DownlinkReq{RespType:"downlink", DownlinkMsg:downlinkMsg}
        devDownlinkCache.Set(cacheId, downlinkPayload, cache.DefaultExpiration)
        w.WriteHeader(http.StatusOK)
    }
}

// Functions users use to make an SSH request to schedule or stop
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
        checkErr(err, "Couldn't insert ssh downlink log")
        device, err := devEnv.db.ADevice(devTable, devName)
        if tunnelStatus == "schedule" {
            if err != nil {
                log.Println("SSH Key not found")
            }
            err = ssh.DelDeviceKey(email_tbl, device.SSHKey)
            port := ssh.AddDeviceKey(email_tbl, device.SSHKey)
            sch := new(scheduler.Sched)
            //TODO: Add device session time
            sch.InitScheduler(time.Hour*1, ssh.DelDeviceKey, email_tbl, device.SSHKey)
            sch.Start()
            sshReq := &DownlinkReq{RespType:tunnelStatus,Port:port}
            devDownlinkCache.Set(cacheId, sshReq, cache.DefaultExpiration)
        }
        if tunnelStatus == "launch" {
        }
        if tunnelStatus == "stop" {
            log.Println("Stopping SSH Tunnel")
            err = ssh.DelDeviceKey(email_tbl, device.SSHKey)
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
        // Return back msg of typle UsrDownlinkReq
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
    if r.Method == "POST" {
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
