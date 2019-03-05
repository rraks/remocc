package controller


import (
    "net/http"
    "github.com/satori/go.uuid"
    "github.com/patrickmn/go-cache"
    "time"
    "github.com/rraks/remocc/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var usrAuthCache *cache.Cache

func init() {
    usrAuthCache = cache.New(2*time.Hour,4*time.Hour)
}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
        }
        sessionEmail := emailCook.Value
        sessionToken := tokenCook.Value
        if val, found := usrAuthCache.Get(sessionEmail); found {
            if sessionToken == val {
                fn(w, r)
                return
            }
        }
        http.Redirect(w, r, "/login/", http.StatusFound)
    }
}


func LogSession(w http.ResponseWriter, usr *models.User) {
    sessionToken := uuid.NewV4().String()
    usrAuthCache.Set(usr.Email, sessionToken, cache.DefaultExpiration)
    http.SetCookie(w, &http.Cookie{
        Name:    "dev_table",
        Value:   "devices_" + usr.Name,
        Path: "/",
    })
    http.SetCookie(w, &http.Cookie{
        Name:    "email",
        Value:   usr.Email,
        Path: "/",
    })
    http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
        Value:   sessionToken,
        Path: "/",
    })
}
