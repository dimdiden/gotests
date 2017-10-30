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
func GetWeather(city, country string) Weather {
	url := ROOT_URL + API_KEY + "&q=" + city + "," + country
	var weather Weather

	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(body, &weather); err != nil {
		panic(err)
	}
	return weather
}

func (w *Weather) createBulk(days int) (bulk [][]string) {
	var countday int
	var rowcounter int
	current_day := strings.Split(w.List[0].Timestamp, " ")[0] // first day in the list

	for _, msrmnt := range w.List {
		celsius := FartoCel(msrmnt.Values.Temp)
		datetime := strings.Split(msrmnt.Timestamp, " ") // []string{date, time}

		if current_day != datetime[0] {
			countday++
			current_day = datetime[0]
		}
		if countday == days {
			return bulk[:rowcounter]
		}

		row := append(datetime, strconv.FormatFloat(celsius, 'f', -1, 32), strconv.Itoa(msrmnt.Values.Hum))
		bulk = append(bulk, row)
		rowcounter++
	}
	return bulk[:rowcounter]
}

func (w *Weather) Render(days int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"date", "time", "temperature", "humidity"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("|")
	table.SetCaption(true, w.City.Name)

	bulk := w.createBulk(days)
	table.AppendBulk(bulk)

	table.Render()
}

func FartoCel(f float64) float64 {
	return math.Ceil(f - 273.15)
}

// =============================================================
func main() {

	city := flag.String("city", "Kyiv", "Choose the target city")
	country := flag.String("country", "ua", "Choose the country")
	days := flag.Int("d", 1, "Number of the displayed days")
	flag.Parse()

	w := GetWeather(*city, *country)
	w.Render(*days)
}
