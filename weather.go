package main

import (
  "os"
  // "fmt"
  "net/http"
  "io/ioutil"
  "log"
  "encoding/json"
  "math"
  "github.com/olekukonko/tablewriter"
  "strconv"
)

const (
  ROOT_URL = "http://api.openweathermap.org/data/2.5/forecast?appid="
  API_KEY = "dcf5b77beaf67157ac55a0263f8def87"
)

type Weather struct {
  City *City
  List []*Measurement
}

type City struct {
  Name string
}

type Measurement struct {
  Date string `json:"dt_txt"`
  Values *Values `json:"main"`
}

type Values struct {
  Temp float64 `json:"temp"`
  Hum int `json:"humidity"`
}

func (v *Values) Stringify() []string {
  celsius := fartoCel(v.Temp)
  return []string{strconv.FormatFloat(celsius, 'f', -1, 32), strconv.Itoa(v.Hum)}
}

func (m *Measurement) Stringify() []string {
  s := []string{m.Date}
  return append(s, m.Values.Stringify()...)
}

// http://api.openweathermap.org/data/2.5/forecast?appid=dcf5b77beaf67157ac55a0263f8def87&q=Sumy,ua
func getData() []byte {
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

func fartoCel(f float64) float64  {
  return math.Ceil(f - 273.15)
}

// ================================
func main() {
  data := getData()

  var w Weather
  if err := json.Unmarshal(data, &w); err != nil {
    panic(err)
  }

  table := tablewriter.NewWriter(os.Stdout)
  table.SetHeader([]string{"time", "temperature", "humidity"})

  for _, msrmnt := range w.List {
    table.Append(msrmnt.Stringify())
    // fmt.Println(msrmnt.Date, "||", fartoCel(msrmnt.Values.Temp))
  }

  table.Render()
}
