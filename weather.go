package main

import (
	"os"
	// "fmt"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const (
	ROOT_URL = "http://api.openweathermap.org/data/2.5/forecast?appid="
	API_KEY  = "dcf5b77beaf67157ac55a0263f8def87"
)

type Weather struct {
	City struct {
		Name string
	}
	List []struct {
		Timestamp string `json:"dt_txt"`
		Values    struct {
			Temp float64
			Hum  int `json:"humidity"`
		} `json:"main"`
	}
}

// http://api.openweathermap.org/data/2.5/forecast?appid=dcf5b77beaf67157ac55a0263f8def87&q=Sumy,ua
func (w *Weather) GetData() {
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

	if err := json.Unmarshal(body, &w); err != nil {
		panic(err)
	}
}

func (w *Weather) createBulk() [][]string {
	var data [][]string
	for _, msrmnt := range w.List {
		celsius := FartoCel(msrmnt.Values.Temp)
		datetime := strings.Split(msrmnt.Timestamp, " ") // []string{date, time}
		row := append(datetime, strconv.FormatFloat(celsius, 'f', -1, 32), strconv.Itoa(msrmnt.Values.Hum))
		data = append(data, row)
	}
	return data
}

func (w *Weather) Render() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"date", "time", "temperature", "humidity"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("|")
	table.SetCaption(true, w.City.Name)

	bulk := w.createBulk()
	table.AppendBulk(bulk)

	table.Render()
}

func FartoCel(f float64) float64 {
	return math.Ceil(f - 273.15)
}

// ================================
func main() {
	// data := getData()

	var w Weather
	// if err := json.Unmarshal(data, &w); err != nil {
	// 	panic(err)
	// }
	w.GetData()
	w.Render()
}
