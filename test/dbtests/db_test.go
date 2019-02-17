package tests

import (
    "testing"
    "remocc/pkg/models"
)

type Env struct {
    db models.Datastore
}

type mockDB struct{}


func TestOpenDB(t *testing.T)  {
    db, err := models.InitDB()
    if err != nil {
        t.Errorf("Failed to open DB")
        panic(err)
    }
    t.Log("DB Opened successfully")
    _ = db
}


func TestGetUsers(t *testing.T) {
    db, err := models.InitDB()
    if err != nil {
        t.Errorf("Failed to open DB")
        panic(err)
    }
    t.Log("DB Opened successfully")
    env := &Env{db}
    usrs, err := env.db.AllUsers()
    if err != nil {
        t.Errorf(err.Error())
    }
    for _, usr := range usrs {
        t.Log(usr.Name)
    }
}


