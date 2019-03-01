package controller


import (
    "net/http"
    "github.com/patrickmn/go-cache"
    "time"
    "log"
)

var usrAuthCache *cache.Cache

func init() {
    usrAuthCache = cache.New(2*time.Hour,4*time.Hour)
}

func ProvideHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenCook, err1 := r.Cookie("session_token")
        emailCook, err2 := r.Cookie("email")
        if err1 != nil || err2 !=nil {
            if err1 == http.ErrNoCookie || err2 == http.ErrNoCookie {
                http.Redirect(w, r, "/login/", http.StatusFound)
                return
            }
            sessionEmail := emailCook.Value
            sessionToken := tokenCook.Value
            log.Println(sessionEmail)
            if val, found := usrAuthCache.Get(sessionEmail); found {
                if sessionToken == val {
                    http.ServeFile(w,r, "web/views/base.html")
                    return
                }
            }
            http.Redirect(w, r, "/login/", http.StatusFound)
            return
        }
        fn(w, r)
    }
}

