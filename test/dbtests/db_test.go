package tests

import (
    "testing"
    "github.com/rraks/remocc/pkg/models"
)

type Env struct {
    db models.Datastore
}

type mockDB struct{}

func openDB()(db *models.DB)  {
    db, err := models.InitDB()
    if err != nil {
        panic(err)
    }
    return db
}




func TestDBMain(t *testing.T) {

    testTables := make([]models.User,8)
    testTables =  []models.User{{"Rax", "1ax@abc.com", "abc", "xyz"},
                                 {"qax", "2ax@abc.com", "abc", "xyz"},
                                 {"wax", "3ax@abc.com", "abc", "xyz"},
                                 {"eax", "4ax@abc.com", "abc", "xyz"},
                                 {"tax", "5ax@abc.com", "abc", "xyz"},
                                 {"aax", "6ax@abc.com", "abc", "xyz"},
                                 {"cax", "7ax@abc.com", "abc", "xyz"},
                                 {"vax", "8ax@abc.com", "abc", "xyz"}, }

    db := openDB()
    env := &Env{db}
    t.Log("DB opened succesffully")


    // Delete group
    err := env.db.DeleteGrp("xyz")
    if err != nil {
        t.Errorf(err.Error())
    }
    t.Log("Deleted Group", "xyz")


    for i := 0; i<8; i++ {
        //Create test Users
        id, err := env.db.NewUser(testTables[i].Name, testTables[i].Email, 
                    testTables[i].Org, testTables[i].Grp)
        if err != nil {
            t.Errorf(err.Error())
        }
        t.Log("Created user", testTables[i].Name, "\t id ", id)
    }

    // Get one User
    usr, err := env.db.AUser("wax")
    if err != nil {
        t.Errorf(err.Error())
    }
    t.Log("Found user", usr.Name)

    //Delete that user
    err = env.db.DeleteUser("wax")
    if err != nil {
        t.Errorf(err.Error())
    }
    t.Log("Deleted user", "wax")

    // Get all Users
    usrs, err := env.db.AllUsers()
    if err != nil {
        t.Errorf(err.Error())
    }
    for _, usr := range usrs {
        t.Log("Found user  ", usr.Name)
    }

    // Delete group again
    err = env.db.DeleteGrp("xyz")
    if err != nil {
        t.Errorf(err.Error())
    }
    t.Log("Deleted Group", "xyz")

}
