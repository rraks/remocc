package models



type User struct {
    Name, Email, Org, Grp string
}


func (db *DB) NewUser(name string, email string, org string, grp string) (int, error) {
    query := "INSERT INTO users (name, email, orgname, groupname) VALUES ($1, $2, $3, $4) RETURNING id"
    id := 0
    err := db.QueryRow(query, name, email, org, grp).Scan(&id)
    if err != nil {
        return id, err
    }
    return id, nil
}


func (db *DB) AllUsers() ([]*User, error) {
    usrs := make([]*User, 0)
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    id := 0
    for rows.Next() {
        usr := new(User)
        err := rows.Scan(&id, &usr.Name,&usr.Email,&usr.Org,&usr.Grp)
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


func (db *DB) AUser(name string) (*User, error) {
    usr := new(User)
    rows, err := db.Query("SELECT * FROM users WHERE Name='" + name + "'")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    id := 0
    rows.Next()
    err = rows.Scan(&id, &usr.Name,&usr.Email,&usr.Org,&usr.Grp)
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



