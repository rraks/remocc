package controller

import (
    "net/http"
    "log"
    "github.com/rraks/remocc/pkg/models"
	"golang.org/x/crypto/bcrypt"
    "github.com/satori/go.uuid"
    "github.com/patrickmn/go-cache"
    "hash/fnv"
    "bytes"
    "time"
)


type Env struct {
    db models.UserStore
}

var env *Env
var authCache *cache.Cache

func init() {
    db, err := models.InitDB()
    if err != nil {
        panic(err)
    }
    env = &Env{db}

    authCache = cache.New(2*time.Hour,4*time.Hour)
}


func hashString(s string) string  {
    var out []byte
    h := fnv.New32a()
    log.Println("Received", s)
    h.Write([]byte(s))
    h.Sum(out)
    log.Println(out)
    n := bytes.IndexByte(out, 0)
    log.Println("Size ", n)
    return string(out[:n])
}



func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    tokenCook, err1 := r.Cookie("session_token")
    emailCook, err2 := r.Cookie("email")
    if err1 != nil || err2 !=nil {
        if err1 == http.ErrNoCookie || err2 == http.ErrNoCookie {
            // If the cookie is not set, return an unauthorized status
            http.Redirect(w, r, "/login/", http.StatusFound)
            return
        }
        http.Redirect(w, r, "/login/", http.StatusFound)
        return
    }
    sessionEmail := emailCook.Value
    sessionToken := tokenCook.Value
    if val, found := authCache.Get(sessionEmail); found {
        if sessionToken == val {
            authCache.Delete(sessionEmail)
        }
    }
    http.Redirect(w, r, "/login/", http.StatusFound)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        http.ServeFile(w,r, "web/views/login.html")
    } else if r.Method == "POST" {
        err := r.ParseForm()
        if err != nil {
            log.Println("Failed to parse form")
        }
        uname := r.Form["email"][0]
        usr, err := env.db.AUser(uname)
        if err != nil {
            http.Redirect(w, r, "/login", http.StatusFound) // TODO: Add not found flash message
            return
        }
        email := usr.Email
        hash, err := env.db.GetPwd(email)
        match := CheckPasswordHash(r.Form["password"][0], hash)
        if match == true {
            sessionToken := uuid.NewV4().String()
            authCache.Set(email, sessionToken, cache.DefaultExpiration)
            http.SetCookie(w, &http.Cookie{
                    Name:    "email",
                    Value:   email,
                    Path: "/",
                })
            http.SetCookie(w, &http.Cookie{
                    Name:    "session_token",
                    Value:   sessionToken,
                    Path: "/",
                })
            http.Redirect(w,r, "/", http.StatusFound)
            return
        }
        http.Redirect(w,r, "/login/", http.StatusFound)
    }
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        http.ServeFile(w,r, "web/views/register.html")
    } else if r.Method == "POST" {
        err := r.ParseForm()
        if err != nil {
            http.Redirect(w, r, "/register/", http.StatusFound)
        }
        name := r.Form["name"][0]
        email := r.Form["email"][0]
        pswd := r.Form["password"][0]
        cnfmpswd := r.Form["confPassword"][0]
        if  (pswd == cnfmpswd) {
            hashpwd, _ := HashPassword(pswd)
            err1 := env.db.CreateDeviceTable("dev_"+name) // TODO: Replace with hash of email instead of name
            err2 := env.db.CreateAppTable("app_"+name)
            if (err1 != nil) && (err2 != nil) {
                http.Redirect(w, r, "/register/", http.StatusFound)
                return
            }
            _, err := env.db.NewUser(name, email, "default", "default", hashpwd, "dev"+name, "app"+name)
            if err != nil {
                http.Redirect(w, r, "/register/", http.StatusFound)
                return 
            }
            http.Redirect(w, r, "/login/", http.StatusFound)
            return
        }
        http.Redirect(w, r, "/register/", http.StatusFound)
    }
}

func FrontPageHandler(w http.ResponseWriter, r *http.Request) {
    tokenCook, err1 := r.Cookie("session_token")
    emailCook, err2 := r.Cookie("email")
    if err1 != nil || err2 !=nil {
        if err1 == http.ErrNoCookie || err2 == http.ErrNoCookie {
            http.Redirect(w, r, "/login/", http.StatusFound)
            return
        }
        http.Redirect(w, r, "/login/", http.StatusFound)
        return
    }
    sessionEmail := emailCook.Value
    sessionToken := tokenCook.Value
    if val, found := authCache.Get(sessionEmail); found {
        if sessionToken == val {
            http.ServeFile(w,r, "web/views/base.html")
            return
        }
    }
    http.Redirect(w, r, "/login/", http.StatusFound)
}

