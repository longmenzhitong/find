// package weather implements methods for handling the 'weather' order.
package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// key of FIND in amap
const key = "205490cab20372f34b57fbdadf28de90"

type weather struct {
	Status   string
	Count    string
	Info     string
	Infocode string
	Lives    []live
	Forecast forecast
}

type live struct {
	Province      string
	City          string
	Adcode        string
	Weather       string
	Temperature   string
	Winddirection string
	Windpower     string
	Humidity      string
	Reporttime    string
}

type forecast struct {
	City       string
	Adcode     string
	Province   string
	Reporttime string
	Casts      []cast
}

type cast struct {
	Date         string
	Week         string
	Dayweather   string
	Nightweather string
	Daytemp      string
	Nighttemp    string
	Daywind      string
	Nightwind    string
	Daypower     string
	Nightpower   string
}

type geo struct {
	Status   string
	Count    string
	Info     string
	Geocodes []geocode
}

type geocode struct {
	Adcode string
}

// Search is used to query weather from amap and print it.
func Search(address string, all bool) error {
	mod := "base"
	if all {
		mod = "all"
	}
	adcode, err := getAdcode(address)
	if err != nil {
		return fmt.Errorf("get adcode error: %v", err)
	}
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%s&extensions=%s", key, adcode, mod)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get weather error: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read weather error: %v", err)
	}

	var _weather weather
	err = json.Unmarshal(body, &_weather)
	if err != nil {
		return fmt.Errorf("parse weather error: %v", err)
	}

	_live := _weather.Lives[0]
	fmt.Printf("Province：%s\n", _live.Province)
	fmt.Printf("City：%s\n", _live.City)
	fmt.Printf("Weather: %s\n", _live.Weather)
	fmt.Printf("Temperature: %s\n", _live.Temperature)
	fmt.Printf("Wind Direction: %s\n", _live.Winddirection)
	fmt.Printf("Wind Power: %s\n", _live.Windpower)
	fmt.Printf("Humidity: %s\n", _live.Humidity)
	fmt.Printf("Report Time: %s\n", _live.Reporttime)

	return nil
}

// getAdcode is used to get amap's adcode for searching weather,
// returning amap's adcode and error.
func getAdcode(address string) (string, error) {
	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?key=%s&address=%s&output=JSON", key, address)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("get geo error: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read geo error: %v", err)
	}

	var _geo geo
	err = json.Unmarshal(body, &_geo)
	if err != nil {
		return "", fmt.Errorf("parse geo error: %v", err)
	}
	if len(_geo.Geocodes) > 1 {
		return "", fmt.Errorf("multiple geocode error")
	}
	return _geo.Geocodes[0].Adcode, nil
}
