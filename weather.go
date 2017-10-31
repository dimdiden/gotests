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
	List []Measurement
}

type Measurement struct {
	Timestamp string `json:"dt_txt"`
	Values    struct {
		Temp float64
		Hum  int     `json:"humidity"`
		Pres float64 `json:"pressure"`
	} `json:"main"`
	Description []struct {
		Main     string
		Specific string `json:"description"`
	} `json:"weather"`
	Clouds struct {
		All int
	}
}

// http://api.openweathermap.org/data/2.5/forecast?appid=dcf5b77beaf67157ac55a0263f8def87&q=Sumy,ua
func GetWeather(city, country string) (w Weather) {
	url := ROOT_URL + API_KEY + "&q=" + city + "," + country

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
	return
}

func (m *Measurement) createRow() (row []string) {
	daytime := strings.Split(m.Timestamp, " ") // []string{date, time}

	day, time := daytime[0], daytime[1]
	celsius := strconv.FormatFloat(FartoCel(m.Values.Temp), 'f', -1, 32)
	humidity := strconv.Itoa(m.Values.Hum)
	pressure := strconv.FormatFloat(m.Values.Pres, 'f', -1, 32)
	main := m.Description[0].Main
	specific := m.Description[0].Specific
	cloud := strconv.Itoa(m.Clouds.All)

	row = append(row, day, time, celsius, main, specific, humidity, pressure, cloud)
	return
}

func (w *Weather) createBulk(days int, ncol int) (bulk [][]string) {
	var countday int
	var rowcounter int
	current_day := strings.Split(w.List[0].Timestamp, " ")[0] // first day in the list "2017-10-31"

	for _, m := range w.List {
		day := strings.Split(m.Timestamp, " ")[0]

		if current_day != day {
			countday++
			current_day = day
		}
		if countday == days {
			return bulk[:rowcounter]
		}

		row := m.createRow()[:ncol]
		bulk = append(bulk, row)
		rowcounter++
	}
	return bulk[:rowcounter]
}

func (w *Weather) Render(days int, ncol int) {
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"date", "time", "temperature", "main", "specific", "humidity", "pressure", "cloud"}
	table.SetHeader(header[:ncol])
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("|")

	caption := "City: " + w.City.Name + "; displayed days: " + strconv.Itoa(days)
	table.SetCaption(true, caption)

	bulk := w.createBulk(days, ncol)
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
	full := flag.Bool("f", false, "Enables detailed description")
	flag.Parse()

	var ncol int // number of columns
	if *full {
		ncol = 8
	} else {
		ncol = 5
	}

	w := GetWeather(*city, *country)
	w.Render(*days, ncol)
}
