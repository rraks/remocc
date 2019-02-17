package models



type User struct {
    Name, Email, Org, Grp string
}


func (db *DB) AllUsers() ([]*User, error) {
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    usrs := make([]*User, 0)
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

