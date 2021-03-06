package controller

import (
    "net/http"
    "log"

)

//Server connection parameters
const (
    host = "localhost"
    port = 5600
)


func Start() {
    //mux for path parameterd api endpoints 
    mux := http.NewServeMux()

    //Serve static resources
    fs := http.FileServer(http.Dir("web/static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    //User Login logout
    mux.HandleFunc("/login/", LoginHandler)
    mux.HandleFunc("/logout/", LogoutHandler)

    //Main page Router
    mux.HandleFunc("/register/", RegisterHandler)

    //Test Handlers
    //mux.HandleFunc("/", Testprovidehandler(FrontPageHandler))
    mux.HandleFunc("/user/devices/info/test", Testprovidehandler(DeviceLogsHandler))
    //mux.HandleFunc("/user/devices/downlink/", Testprovidehandler(SendDeviceDownlink))
    //mux.HandleFunc("/user/devices/ssh/", Testprovidehandler(SendSSHRequest))

    //Actual
    mux.HandleFunc("/", ProvideWebHandler(FrontPageHandler))
    mux.HandleFunc("/user/devices/info/", ProvideWebHandler(DeviceInfoHandler))
    mux.HandleFunc("/user/devices/info/ssh/", ProvideWebHandler(DeviceSSHStatusHandler))
    mux.HandleFunc("/user/devices/downlink/", ProvideWebHandler(DeviceDownlinkHandler))
    mux.HandleFunc("/user/devices/ssh/", ProvideWebHandler(SSHRequestHandler))
    mux.HandleFunc("/user/devices/manage/", ProvideWebHandler(DeviceManagerHandler))

    //Device Handlers
    mux.HandleFunc("/devices/login/", DeviceLoginHandler)

    mux.HandleFunc("/devices/data/heartbeat/", ProvideApiHandler(HeartBeatHandler))
    mux.HandleFunc("/devices/data/uplink/", ProvideApiHandler(DeviceUplinkHandler))
    mux.HandleFunc("/devices/data/ssh/", ProvideApiHandler(DeviceTunnelRequestHandler))


    //Serve 
    log.Fatal(http.ListenAndServe(":3000", mux))
}
