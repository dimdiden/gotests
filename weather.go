package main

import (
	"encoding/json"
	"flag"
	// "fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
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
func (w *Weather) GetData(city, country *string) {
	url := ROOT_URL + API_KEY + "&q=" + *city + "," + *country

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

// =============================================================
func main() {
	var w Weather

	city := flag.String("city", "Sumy", "Choose the target city")
	country := flag.String("country", "ua", "Choose the country")
	flag.Parse()

	w.GetData(city, country)
	w.Render()
}
