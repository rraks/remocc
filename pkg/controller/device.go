package controller

//TODO :[
//          Verify uniqueness of name
//]

import (
    "net/http"
    "encoding/json"
    "github.com/rraks/remocc/pkg/models"
    "github.com/rraks/remocc/pkg/views"
    "log"
    "github.com/dgrijalva/jwt-go"
    "errors"
    "github.com/mitchellh/mapstructure"
)


type DevEnv struct {
    db models.DeviceStore
}

var devEnv *DevEnv


type DevReq struct {
    DevName string `json: "devName"`
    UName string `json: "uName"`
    Pwd string `json: "pwd"`
}

type JWToken struct {
    Token string `json:"token"`
}

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
        err = devEnv.db.CreateDeviceLog(devName) // TODO: Have a better way of creating device log
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
    views.RenderTableRow(w, allDevices) // Ignore error
}


// TODO: Make this agnostic to the user
func DeviceLoginHandler(w http.ResponseWriter, r *http.Request) {
    var devReq DevReq
    json.NewDecoder(r.Body).Decode(&devReq)
    hash, err := devEnv.db.GetDevPwd("devices_"+devReq.UName, devReq.DevName)
    if err != nil {
        log.Println(err)
    }
    match := CheckPasswordHash(devReq.Pwd, hash)
    // Create token, TODO :check user policies
    if match == true {
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "uName":  devReq.UName,
            "devName":  devReq.DevName,
            "pwd":  devReq.Pwd,
        })
        tokenString, err := token.SignedString([]byte("password")) // TODO : replace in production through init 
        if err != nil {
            log.Println(err)
        }
        json.NewEncoder(w).Encode(JWToken{Token: tokenString})
    }
}

func DeviceDataHandler(w http.ResponseWriter, r *http.Request) {
    key := r.Header.Get("authToken")
    token, _ := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("Failed to validate token")
        }
        return []byte("password"), nil // TODO : replace in production through init
    })
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        var devReq DevReq
        mapstructure.Decode(claims, &devReq)
        log.Println(devReq)
        log.Println("Starting db check")
        passwordHash, err := devEnv.db.GetDevPwd("devices_"+devReq.UName, devReq.DevName)
        if err != nil {
            log.Println("Not found")
        }
        log.Println("Stopping db check")
        if ok = CheckPasswordHash(devReq.Pwd, passwordHash); ok {
            log.Println("Success")
        }

    }
}



