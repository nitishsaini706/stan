package main

import (
    "fmt"
    "log"
    "sync"
    "github.com/nitishsaini706/stan/assignment2/models"
    "github.com/nitishsaini706/stan/assignment2/store"
    "github.com/nitishsaini706/stan/assignment2/web"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Initialize SQLite database
    db, err := sql.Open("sqlite3", "./user.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Migrate or ensure tables are created
    err = store.Migrate(db)
    if err != nil {
        log.Fatal("Failed to migrate or create tables:", err)
    }

    s := store.New(db)
    
    // Concurrently populate the store with initial data
    wg := sync.WaitGroup{}
    users := []models.User{
        {ID: 1, Name: "Nitish Saini", Email: "nitish@saini.com"},
        {ID: 2, Name: "Nitish 2", Email: "nitish@2.com"},
    }
	
    for _, user := range users {
        wg.Add(1)
        go func(u models.User) {
            defer wg.Done()
            if err := s.CreateUser(u); err != nil {
                fmt.Printf("Error creating user %v: %v\n", u, err)
            }
        }(user)
    }
    wg.Wait()

    // Setup and start the web server
    r := web.SetupRouter(s)
    log.Fatal(r.Run(":8080"))
}

