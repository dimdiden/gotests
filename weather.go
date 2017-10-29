package main

import (
  // "os"
  "fmt"
  "net/http"
  "io/ioutil"
  "log"
  "encoding/json"
)

const (
  ROOT_URL = "http://api.openweathermap.org/data/2.5/forecast?appid="
  API_KEY = "dcf5b77beaf67157ac55a0263f8def87"
)

type Weather struct {
  Code string `json:"cod"`
  City *City `json:"city"`
}

type City struct {
  Name string `json:"name"`
}

func getData() ([]byte) {
  url := ROOT_URL + API_KEY + "&q=Sumy,ua"

  res, err := http.Get(url)
  if err != nil {
    log.Fatalln(err)
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    panic(err.Error())
  }

  return body
}

// http://api.openweathermap.org/data/2.5/forecast?appid=dcf5b77beaf67157ac55a0263f8def87&q=Sumy,ua
func main() {
  body := getData()

  w := &Weather{}
  if err := json.Unmarshal(body, &w); err != nil {
    panic(err)
  }

  fmt.Println(w.Code, w.City.Name)
}
