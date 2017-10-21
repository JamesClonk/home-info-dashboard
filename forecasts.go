package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}
	return data, nil
}

func parseXML(data []byte) {
	var dat interface{}
	if err := xml.Unmarshal(data, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
}

func forecasts() {
	data, err := getData("http://www.yr.no/place/Switzerland/Bern/Bützberg/forecast.xml")
	if err != nil {
		log.Printf("Get forecast data: %v\n", err)
		return
	}
	parseXML(data)
}
