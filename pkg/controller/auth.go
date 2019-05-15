package controller


// auth.go
// Provides authorization closures and structures for tokens/claims 

import (
    "net/http"
    "github.com/satori/go.uuid"
    "github.com/patrickmn/go-cache"
    "time"
    "github.com/rraks/remocc/pkg/models"
	"golang.org/x/crypto/bcrypt"
    "log"
    "github.com/dgrijalva/jwt-go"
    "errors"
    "encoding/json"
    "github.com/mitchellh/mapstructure"
    "strings"
    "os"
)

// Memory cache for user session tokens
var usrAuthCache *cache.Cache

// JWT signing password
var jwtPassword string

// Claims structure for devices claiming a JWT
type DevClaims struct {
    DevName string `json: "devName"`
    Email string `json: "email"`
    Pwd string `json: "pwd"`
}

// JWT response message structure
type JWToken struct {
    Token string `json:"token"`
}


func init() {
    usrAuthCache = cache.New(2*time.Hour,4*time.Hour)
    jwtPassword = os.Getenv("JWT_PASSWORD")
}

// Simple BCrypt hash to store user and device passwords
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

// Check BCrypt hash by matching
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


// Log a users web session via cookies
func LogSession(w http.ResponseWriter, usr *models.User) {
    sessionToken := uuid.NewV4().String()
    usrAuthCache.Set(usr.Email, sessionToken, cache.DefaultExpiration)
    email_tbl := strings.Replace(usr.Email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    http.SetCookie(w, &http.Cookie{
        Name:    "dev_table",
        Value:   "devices_" + email_tbl,
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



// A test closure to  wrap web functions for testing. A user "a" with email "a@a.com" password "a" is used
func Testprovidehandler(fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
                fn(w, r, "a@a.com", "devices_a_a_com")
    }
}

// Closure to provide web content. Checks for cookies and returns to the requester the requested page
// if user's session token is present in usrAuthCache. Also provides to the API function the users email and the device table name
// TODO: No need of device table. It can be constructed from user email
func ProvideWebHandler(fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenCook, err1 := r.Cookie("session_token")
        emailCook, err2 := r.Cookie("email")
        devTableCook, err3 := r.Cookie("dev_table")
        if err1 != nil || err2 !=nil || err3 != nil {
            if err1 == http.ErrNoCookie || err2 == http.ErrNoCookie {
                http.Redirect(w, r, "/login/", http.StatusFound)
                return
            }
        }
        sessionEmail := emailCook.Value
        sessionToken := tokenCook.Value
        sessionDevTable := devTableCook.Value
        if val, found := usrAuthCache.Get(sessionEmail); found {
            if sessionToken == val {
                fn(w, r, sessionEmail, sessionDevTable)
                return
            }
        }
        http.Redirect(w, r, "/login/", http.StatusFound)
    }
}



// Provides JWT token when presented with valid device claims. Checks device name, user's email, and devices password
// to return back a JWT 
// TODO: Make this agnostic to the user
// TODO: Add token validity period
func DeviceLoginHandler(w http.ResponseWriter, r *http.Request) {
    var devClaims DevClaims
    json.NewDecoder(r.Body).Decode(&devClaims)
    email_tbl := strings.Replace(devClaims.Email,"@","_",-1)
    email_tbl = strings.Replace(email_tbl,".","_",-1)
    hash, err := devEnv.db.GetDevPwd("devices_"+email_tbl, devClaims.DevName)
    if err != nil {
        log.Println(err)
    }
    match := CheckPasswordHash(devClaims.Pwd, hash)
    // Create token, TODO :check user policies
    if match == true {
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "email":  devClaims.Email,
            "devName":  devClaims.DevName,
            "pwd":  devClaims.Pwd,
        })
        tokenString, err := token.SignedString([]byte(jwtPassword)) // TODO : replace in production through init 
        if err != nil {
            log.Println(err)
        }
        json.NewEncoder(w).Encode(JWToken{Token: tokenString})
    }
}


// Closure to provide requested API for devices. Checks for JWT token claims and services the requested api if
// token is valid. In addition, provides the devices claim to the underlying API
func ProvideApiHandler(fn func(http.ResponseWriter, *http.Request, *DevClaims)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        key := r.Header.Get("authToken")
        token, _ := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("Failed to validate token")
            }
            return []byte(jwtPassword), nil // TODO : replace in production through init
        })
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            var devClaims DevClaims
            mapstructure.Decode(claims, &devClaims)
            email_tbl := strings.Replace(devClaims.Email,"@","_",-1)
            email_tbl = strings.Replace(email_tbl,".","_",-1)
            passwordHash, err := devEnv.db.GetDevPwd("devices_"+email_tbl, devClaims.DevName)
            if err != nil {
                w.Write([]byte("Invalid authorization"))
            }
            if ok = CheckPasswordHash(devClaims.Pwd, passwordHash); ok {
                fn(w, r, &devClaims)
                return
            }
        }
        w.Write([]byte("Invalid authorization"))
    }
}
