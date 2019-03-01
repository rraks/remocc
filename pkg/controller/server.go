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
    // Serve static resources
    fs := http.FileServer(http.Dir("web/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.HandleFunc("/login/", LoginHandler)
    http.HandleFunc("/logout/", LogoutHandler)
    http.HandleFunc("/", ProvideHandler(FrontPageHandler))
    http.HandleFunc("/register/", RegisterHandler)
    log.Fatal(http.ListenAndServe(":3000", nil))
}
