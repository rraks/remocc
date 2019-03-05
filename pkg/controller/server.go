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
    //http.HandleFunc("/", ProvideHandler(FrontPageHandler))
    mux.HandleFunc("/", FrontPageHandler)
    mux.HandleFunc("/register/", RegisterHandler)

    // Device Handlers
    //http.HandleFunc("/device/", ProvideHandler(DeviceManagerHandler))
    mux.HandleFunc("/devices/", DeviceManagerHandler)
    mux.HandleFunc("/devices/login/", DeviceLoginHandler)
    mux.HandleFunc("/devices/data/", DeviceDataHandler)



    // Serve 
    log.Fatal(http.ListenAndServe(":3000", mux))
}
