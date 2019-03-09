package controller

import (
    "net/http"
    "log"

)

// Server connection parameters
const (
    host = "localhost"
    port = 5600
)


func Start() {
    // mux for path parameterd api endpoints 
    mux := http.NewServeMux()

    // Serve static resources
    fs := http.FileServer(http.Dir("web/static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    //User Login logout
    mux.HandleFunc("/login/", LoginHandler)
    mux.HandleFunc("/logout/", LogoutHandler)

    // Main page Router
    //http.HandleFunc("/", ProvideWebHandler(FrontPageHandler))
    mux.HandleFunc("/register/", RegisterHandler)

    //mux.HandleFunc("/", ProvideWebHandler(FrontPageHandler))
    //mux.HandleFunc("/user/devices/data/", ProvideWebHandler(UserDataHandler))
    mux.HandleFunc("/", Testprovidehandler(FrontPageHandler))
    mux.HandleFunc("/user/devices/data/", Testprovidehandler(UserDataHandler))

    // Device Handlers
    mux.HandleFunc("/devices/login/", DeviceLoginHandler)
    mux.HandleFunc("/devices/data/", ProvideApiHandler(DeviceDataHandler))

    mux.HandleFunc("/user/devices/manage/", DeviceManagerHandler)


    // Serve 
    log.Fatal(http.ListenAndServe(":3000", mux))
}
