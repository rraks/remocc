package models


import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
)



const (
    host = "localhost"
    port = 5432
    user = "postgres"
    password = "password"
    dbname = "users"
)


type UserStore interface {
    AllUsers() ([]*User, error)
    AUser(string) (*User, error)
    NewUser(string, string, string, string, string, string, string) (int, error)
    DeleteUser(string) (error)
    DeleteGrp(string) (error)
    GetPwd(string) (string, error)
    CreateAppTable(string) (error)
    CreateDeviceTable(string) (error)
    DeleteTable(string) (error)
}


type DeviceStore interface {
    AllDevices(string) ([]*Device, error)
    ADevice(string, string) (*Device, error)
    NewDevice(string, string, string, string) (int, error)
    DeleteDevice(string, string) (error)
}

type DB struct {
    *sql.DB
}


func InitDB() (*DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return &DB{db}, nil
}
