package models



type Device struct {
    DevName, MacId, DevDescr string
}


func (db *DB) NewDevice(tableName string, devName string, macId string, devDescr string, sshKey string) (int, error) {
    query := "INSERT INTO "+ tableName + " (devName, macId, devDescr, sshKey) VALUES ($1, $2, $3, $4) RETURNING id"
    id := 0
    err := db.QueryRow(query, devName, macId, devDescr, sshKey).Scan(&id)
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

func (db *DB) DeleteDevice(tableName string, devName string) (error) {
    _, err := db.Exec("DELETE FROM "+ tableName + " where devName=$1",devName)
    if err != nil {
        return err
    }
    return nil
}


