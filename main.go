package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type TemplateJSON struct {
	Status struct {
		Water int `json:"water"`
		Wind int  `json:"wind"`
	} `json:"status"`
}

type WaterWind struct {
	Water int
	WaterStatus string
	Wind int
	WindStatus string
}

var ( 
	PORT = ":8080"
	MIN = 1
	MAX = 100
)

func main() {
	go updateJSONFile()	
	http.HandleFunc("/", getHTML)
	
	fmt.Println("Application is listening on port", PORT)
	
	http.ListenAndServe(PORT, nil)
}

func updateJSONFile() {
	for {
		water := rand.Intn(MAX-MIN) + MIN
		wind := rand.Intn(MAX-MIN) + MIN
	
		data := TemplateJSON{}
		data.Status.Water = water
		data.Status.Wind = wind
		
		bodyFileUpdate, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("error encode struct to json: %v", err)
			return
		}

		err = ioutil.WriteFile("data.json", bodyFileUpdate, 0644)
		if err != nil {
			fmt.Printf("unable to write file: %v", err)
		}

		time.Sleep(15 * time.Second)
	}
}

func getHTML(w http.ResponseWriter, r *http.Request) {
	bodyFile, err := ioutil.ReadFile("data.json")
    if err != nil {
        fmt.Printf("unable to read file: %v", err)
    }
	
	var jsonData = []byte(string(bodyFile))
	var data TemplateJSON
	
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Printf("error decode json to struct: %v", err)
	}

	water := data.Status.Water
	var waterStatus string
	wind := data.Status.Wind
	var windStatus string

	if water >= 6 && water <= 8 {
		waterStatus = "Siaga"
	} else if water > 8 {
		waterStatus = "Bahaya"
	} else {
		waterStatus = "Aman"
	}

	if wind >= 7 && wind <= 15 {
		windStatus = "Siaga"
	} else if wind > 15 {
		windStatus = "Bahaya"
	} else {
		windStatus = "Aman"
	}

	dataForHTML := WaterWind{
		Water: water,
		WaterStatus: waterStatus,
		Wind: wind,
		WindStatus: windStatus,
	}
	
	tpl, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Printf("error parse files: %v", err)
		return
	}

	tpl.Execute(w, dataForHTML)

}
