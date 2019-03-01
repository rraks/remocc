package controller

//TODO :[
//          Verify uniqueness of name
//]

import (
    "net/http"
    "github.com/rraks/remocc/pkg/models"
    "github.com/rraks/remocc/pkg/views"
)


type DevEnv struct {
    db models.DeviceStore
}

var devEnv *DevEnv

func init() {
    db, err := models.InitDB()
    if err != nil {
        panic(err)
    }
    devEnv = &DevEnv{db}
}

func NewDeviceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        //Execute template here
    }
    if r.Method == "POST" {
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
        devName := r.Form["devname"][0]
        macId := r.Form["macId"][0]
        description := r.Form["description"][0]
        _, err = devEnv.db.NewDevice(tableName, devName, macId, description)
        if err != nil {
            http.Redirect(w, r, "/device/", http.StatusFound)
        }
    }
}


func FrontPageHandler(w http.ResponseWriter, r *http.Request) {
    //TODO : Replace with a db call and retreive context via cookie in r
    testDevices := make([]models.Device,2)
    testDevices = []models.Device{
                    {"ABC", "12:12:12:12", "Greatest Device on earth"},
                    {"XYZ", "11:11:12:12", "Not Greatest Device on earth"},
                }
    views.RenderTableRow(w, testDevices)
}
