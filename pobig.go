package main

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "log"
  "encoding/json"
)

type Person struct {
  Id int
  Name string
  Age int
}

func GetPersonList() (personlist []Person) {
  db, err := sql.Open("mysql", "ded:ltleirf@tcp(212.237.13.190:3306)/test")
  if err != nil {
    log.Fatal(err)
  }

  rows, err := db.Query("SELECT * FROM worker")
  if err != nil {
    log.Fatal(err)
  }

  for rows.Next() {
    var p Person
    if err := rows.Scan(&p.Id,&p.Name,&p.Age); err != nil {
      log.Fatal(err)
    }
    personlist = append(personlist, p)
  }

  if err := rows.Err(); err != nil {
    log.Fatal(err)
  }
  return
}

func main() {
  personlist := GetPersonList()
  b, err := json.MarshalIndent(personlist, "", "  ")
  if err != nil {
        log.Fatal(err)
  }
  fmt.Println(string(b))
}
