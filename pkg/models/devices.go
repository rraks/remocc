package models

import (
    "time"
    "strconv"
    "database/sql"
)


type Device struct {
    DevName, MacId, DevDescr string
}

type DeviceLog struct {
    LastSeen time.Time
    TunnelStatus  sql.NullBool
    UplinkMsg sql.NullString
    DownlinkMsg sql.NullString
    PingTime sql.NullInt64
}


func (db *DB) NewDevice(tableName string, devName string, macId string, devDescr string, sshKey string, devPwdHash string) (int, error) {
    query := "INSERT INTO "+ tableName + " (devName, macId, devDescr, sshKey, devPwdHash) VALUES ($1, $2, $3, $4, $5) RETURNING id"
    id := 0
    err := db.QueryRow(query, devName, macId, devDescr, sshKey, devPwdHash).Scan(&id)
    if err != nil {
        return id, err
    }
    return id, nil
}


//func (db *DB) AddAuthToDevice(tableName string, macId string, macId string, description string) (int, error) {
//}


func (db *DB) AllDevices(tableName string) ([]*Device, error) {
    devices := make([]*Device, 0)
    rows, err := db.Query("SELECT devName,macId,devDescr FROM " + tableName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        device := new(Device)
        err := rows.Scan(&device.DevName,&device.MacId,&device.DevDescr)
        if err != nil {
            return nil, err
        }
        devices = append(devices, device)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return devices, nil
}


func (db *DB) ADevice(tableName string, devName string) (*Device, error) {
    device := new(Device)
    rows, err := db.Query("SELECT devName,macId,devDescr FROM " + tableName + " where devName=$1", devName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    rows.Next()
    err = rows.Scan(&device.DevName,&device.MacId,&device.DevDescr)
    if err != nil {
        return nil, err
    }
    return device, nil
}

func (db *DB) GetDeviceLogs( device *Device, offset int, limit int) ([]*DeviceLog, error) {
    deviceLogs := make([]*DeviceLog, 0)
    deviceLog := new(DeviceLog)
    rows, err :=    db.Query("SELECT lastSeen,downlinkMsg,uplinkMsg,pingTime,tunnelStatus FROM " + device.DevName + " LIMIT " + strconv.Itoa(10) )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        err := rows.Scan(&deviceLog.LastSeen,&deviceLog.DownlinkMsg,&deviceLog.UplinkMsg, &deviceLog.PingTime, &deviceLog.TunnelStatus)
        if err != nil {
            return nil, err
        }
        deviceLogs = append(deviceLogs, deviceLog)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return deviceLogs, nil
}


func (db *DB) DeleteDevice(tableName string, devName string) (error) {
    _, err := db.Exec("DELETE FROM "+ tableName + " where devName=$1",devName)
    if err != nil {
        return err
    }
    return nil
}

func (db *DB) DropDeviceTable(tableName string) (error) {
    _, err := db.Exec("DROP TABLE "+ tableName)
    if err != nil {
        return err
    }
    return nil
}



func (db *DB) CreateDeviceTable(tableName string) (error) {
    query := "CREATE TABLE "+ tableName + " (lastSeen timestamptz NOT NULL DEFAULT now(), downlinkMsg text, uplinkMsg text, pingTime smallint, tunnelStatus boolean)"
    _, err := db.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

func (db *DB) GetDevPwd(usersTable string, devName string) (string, error) {
    var hash string
    rows, err := db.Query("SELECT devPwdHash FROM " + usersTable + " WHERE devName='" + devName + "'")
    if err != nil {
        return "", err
    }
    defer rows.Close()
    rows.Next()
    err = rows.Scan(&hash)
    if err != nil {
        return "", err
    }
    return hash, nil
}


func (db *DB) InsertDeviceLog(tableName string, uplinkMsg string, pingTime int) (error) {
    query := "INSERT INTO " + tableName + " (uplinkMsg, pingTime) VALUES ($1, $2) "
    _, err := db.Exec(query, uplinkMsg, pingTime)
    if err != nil {
        return err
    }
    return nil
}
