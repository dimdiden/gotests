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

type City struct {
  Code string `json:"cod"`
  // Temp float32
}

// http://api.openweathermap.org/data/2.5/forecast?appid=dcf5b77beaf67157ac55a0263f8def87&q=Sumy,ua
func main() {
  url := ROOT_URL + API_KEY + "&q=Sumy,ua"

  res, err := http.Get(url)
  if err != nil {
    log.Fatalln(err)
  }
  defer res.Body.Close()

  fmt.Println(res.Status)

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    panic(err.Error())
  }

  // var data map[string]interface{}
  // if err := json.Unmarshal(body, &data); err != nil {
  //   panic(err)
  // }

  c := &City{}
  if err := json.Unmarshal(body, &c); err != nil {
    panic(err)
  }

  fmt.Println(c.Code)

  // pritty, _ := json.MarshalIndent(data, "", "    ")

  // os.Stdout.Write(pritty)

  // decoder := json.NewDecoder(res.Body)

  // body, _ := ioutil.ReadAll(resp.Body)
  // fmt.Println(data)
}
