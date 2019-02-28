package models



type Device struct {
    Name, MacId, Description string
}


func (db *DB) NewDevice(tableName string, name string, macId string, description string) (int, error) {
    query := "INSERT INTO $1 (name, macId, description) VALUES ($2, $3, $4) RETURNING id"
    id := 0
    err := db.QueryRow(query, tableName, name, macId, description).Scan(&id)
    if err != nil {
        return id, err
    }
    return id, nil
}


//func (db *DB) AddAuthToDevice(tableName string, macId string, macId string, description string) (int, error) {
//}


func (db *DB) AllDevices(tableName string) ([]*Device, error) {
    devices := make([]*Device, 0)
    rows, err := db.Query("SELECT deviceName,macId,description FROM $1", tableName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        device := new(Device)
        err := rows.Scan(&device.Name,&device.MacId,&device.Description)
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


func (db *DB) ADevice(tableName string, deviceName string) (*Device, error) {
    device := new(Device)
    rows, err := db.Query("SELECT deviceName,macId,description FROM $1 where deviceName=$2",tableName, deviceName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    rows.Next()
    err = rows.Scan(&device.Name,&device.MacId,&device.Description)
    if err != nil {
        return nil, err
    }
    return device, nil
}

func (db *DB) DeleteDevice(tableName string, deviceName string) (error) {
    _, err := db.Exec("DELETE FROM $1 where deviceName=$2",tableName, deviceName)
    if err != nil {
        return err
    }
    return nil
}


