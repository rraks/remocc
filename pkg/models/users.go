package models


//TODO : [Further modularization to shorten redundant code]
//          

import "log"

type User struct {
    Name, Email, Org, Grp string
}



func (db *DB) NewUser(name string, email string, org string, grp string, pswd string, dev string, app string) (int, error) {
    _, err := db.Query("SELECT name,email,orgname,groupname FROM users WHERE Email='" + email + "'")
    if err == nil{
        query := "INSERT INTO users (name, email, orgname, groupname, password, dev, app) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
        id := 0
        err := db.QueryRow(query, name, email, org, grp, pswd, dev, app).Scan(&id)
        if err != nil {
            return id, err
        }
        return id, nil
    }
    return -1,err
}




func (db *DB) AllUsers() ([]*User, error) {
    usrs := make([]*User, 0)
    rows, err := db.Query("SELECT name,email,orgname,groupname FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        usr := new(User)
        err := rows.Scan(&usr.Name,&usr.Email,&usr.Org,&usr.Grp)
        if err != nil {
            return nil, err
        }
        usrs = append(usrs, usr)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return usrs, nil
}


func (db *DB) AUser(email string) (*User, error) {
    usr := new(User)
    rows, err := db.Query("SELECT name,email,orgname,groupname FROM users WHERE Email='" + email + "'")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    rows.Next()
    err = rows.Scan(&usr.Name,&usr.Email,&usr.Org,&usr.Grp)
    if err != nil {
        return nil, err
    }
    return usr, nil
}

func (db *DB) DeleteUser(email string) (error) {
    _, err := db.Exec("DELETE FROM users where email=$1",email)
    if err != nil {
        return err
    }
    return nil
}

func (db *DB) DeleteGrp(grp string) (error) {
    _, err := db.Exec("DELETE  FROM users where groupname=$1",grp)
    if err != nil {
        return err
    }
    return nil
}

func (db *DB) GetPwd(email string) (string, error) {
    var hash string
    rows, err := db.Query("SELECT password FROM users WHERE Email='" + email + "'")
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


func (db *DB) CreateDevicesTable(tableName string) (error) {
    query := "CREATE TABLE "+ tableName + " (id SERIAL PRIMARY KEY,devName TEXT, macId macaddr, devDescr Text, sshKey Text, devPwdHash Text)"
    _, err := db.Exec(query)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (db *DB) CreateAppTable(tableName string) (error) {
    query := "CREATE TABLE "+ tableName + " (id SERIAL PRIMARY KEY,appName TEXT, description Text, authKey Text)"
    _, err := db.Exec(query)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}


func (db *DB) DeleteTable(tableName string) (error) {
    _, err := db.Exec("DROP TABLE "+ tableName )
    if err != nil {
        return err
    }
    return nil
}
