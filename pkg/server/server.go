package server

import (
    "github.com/rraks/remocc/web"
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

    http.HandleFunc("/", web.LoginHandler)
    http.HandleFunc("/front/", web.FrontPageHandler)
    log.Fatal(http.ListenAndServe(":3000", nil))
}
