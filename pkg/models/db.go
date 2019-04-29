package models


import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
    "os"
    "strconv"
)




// Create users table in postgres using
// create table users(id SERIAL, name VARCHAR(255), email TEXT, orgname VARCHAR(255), groupname VARCHAR(255), password VARCHAR(255), dev TEXT, app TEXT)


type UserStore interface {
    AllUsers() (users []*User, err error)
    AUser(email string) (user *User, err error)
    NewUser(uName string, email string, org string, grp string, pswd string, dev string, app string) (id int, err error)
    DeleteUser(email string) (err error)
    DeleteGrp(grp string) (err error)
    GetPwd(email string) (hash string, err error)
    CreateAppTable(tableName string) (err error)
    CreateDevicesTable(tableName string) (err error)
    DeleteTable(tableName string) (err error)
}


type DeviceStore interface {
    AllDevices(tableName string) (devices []*Device, err error)
    ADevice(tableName string, devName string) (device *Device, err error)
    NewDevice(tableName string, devName string, macId string, devDescr string, sshKey string, devPwdHash string) (serial int, err error)
    DeleteDevice(tableName string, devName string) (err error)
    CreateDeviceTable(tableName string) (err error)
    DropDeviceTable(tableName string) (err error)
    GetDevPwd(usersTable string, devName string) (hash string, err error)
    InsertDeviceUplinkLog(tableName string, uplinkMsg string, pingTime int) (err error)
    InsertDeviceDownlinkLog(tableName string, downlinkMsg string, tunnelStatus string) (err error)
    GetDeviceLogs( device *Device, offset int, limit int) ([]*DeviceLog, error) 
}

type DB struct {
    *sql.DB
}


func InitDB() (*DB, error) {
    host := os.Getenv("POSTGRES_HOST")
    port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
    user := os.Getenv("POSTGRES_USER")
    dbname := os.Getenv("POSTGRES_DB")
    password := os.Getenv("POSTGRES_PASSWORD")
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
