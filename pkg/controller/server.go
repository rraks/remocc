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
    mux.HandleFunc("/register/", RegisterHandler)

    //Test Handlers
//    mux.HandleFunc("/", Testprovidehandler(FrontPageHandler))
//    mux.HandleFunc("/user/devices/info/", Testprovidehandler(GetDeviceInfo))
//    mux.HandleFunc("/user/devices/downlink/", Testprovidehandler(SendDeviceDownlink))
//    mux.HandleFunc("/user/devices/ssh/", Testprovidehandler(SendSSHRequest))

    // Actual
    mux.HandleFunc("/", ProvideWebHandler(FrontPageHandler))
    mux.HandleFunc("/user/devices/info/", ProvideWebHandler(GetDeviceInfo))
    mux.HandleFunc("/user/devices/info/ssh/", ProvideWebHandler(GetDeviceSSHStatus))
    mux.HandleFunc("/user/devices/downlink/", ProvideWebHandler(SendDeviceDownlink))
    mux.HandleFunc("/user/devices/ssh/", ProvideWebHandler(SendSSHRequest))

    // Device Handlers
    mux.HandleFunc("/devices/login/", DeviceLoginHandler)

    mux.HandleFunc("/devices/data/heartbeat/", ProvideApiHandler(SendHeartBeat))
    mux.HandleFunc("/devices/data/uplink/", ProvideApiHandler(SendUplink))
    mux.HandleFunc("/devices/data/ssh/", ProvideApiHandler(MakeTunnelRequest))

    mux.HandleFunc("/user/devices/manage/", DeviceManagerHandler)



    // Serve 
    log.Fatal(http.ListenAndServe(":3000", mux))
}
