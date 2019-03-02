package controller

//TODO :[
//          Verify uniqueness of name
//]

import (
    "net/http"
    "github.com/rraks/remocc/pkg/models"
    "github.com/rraks/remocc/pkg/views"
    "log"
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

    if r.Method == "POST" {
        if err != nil{
            http.Redirect(w, r, "/", http.StatusFound) // TODO: Flash error message here
        }
        devName := r.Form["devName"][0]
        macId := r.Form["macId"][0]
        description := r.Form["devDescr"][0]
        sshKey := r.Form["sshKey"][0]
        _, err = devEnv.db.NewDevice(tableName, devName, macId, description, sshKey)
    }
    //TODO : Won't work, forms can't delete, use ajax
    if r.Method == "DELETE" {
        devName := r.URL.Query().Get("devName")
        confDevName := r.URL.Query().Get("confDevName")
        if devName == confDevName {
            err = devEnv.db.DeleteDevice(tableName, devName)
            if err != nil {
                log.Println(err)
            }
        }
    }
    http.Redirect(w, r, "/", http.StatusFound) // TODO: Flash error message here
}


// TODO: Add a more qualified user context here instead of tableName
func FrontPageHandler(w http.ResponseWriter, r *http.Request) {
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
    allDevices, err := devEnv.db.AllDevices(tableName)
    if err != nil {
        log.Println(err)
    }
    views.RenderTableRow(w, allDevices)
}
