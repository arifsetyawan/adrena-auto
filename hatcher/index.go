package hatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/novalagung/gubrak/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var baseURL = "https://adrena.tech"

type ResponseAuth struct {
	Token struct {
		Access string `json:"access"`
	} `json:"token"`
}

type ResponseTimeTable struct {
	TimeTable []ComponentTimeTableStruct `json:"timetable"`
}

type ComponentTimeTableStruct struct {
	CalDate      string `json:"calDate"`
	IsWorkingDay int    `json:"isWorkingDay"`
}

type ResponseLocation struct {
	Position ComponentPosition `json:"position"`
}

type ComponentPosition struct {
	LocationName string `json:"locationName"`
	ProvinceId   int    `json:"provinceId"`
}

func Auth() (token string, err error) {
	var client = &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
	var data ResponseAuth

	values := map[string]string{"username": os.Getenv("USERNAME"), "password": os.Getenv("password")}
	bodyJSON, _ := json.Marshal(values)
	var payload = bytes.NewBuffer(bodyJSON)

	req, err := http.NewRequest("POST", baseURL+"/auth/api/mobile/session", payload)
	if err != nil {
		return "nil", err
	}

	req.Header.Set("User-Agent", "Adrena%20HR/1 CFNetwork/1325.0.1 Darwin/21.1.0")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	token = data.Token.Access

	return token, fmt.Errorf("auth failed")

}

func IsTodayWorkingOrNot(authToken string) bool {
	var client = &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
	var data ResponseTimeTable

	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")
	req, err := http.NewRequest("GET", baseURL+"/ess/api/timetable/published/week?date="+today, nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", "Adrena%20HR/1 CFNetwork/1325.0.1 Darwin/21.1.0")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	result := gubrak.From(data).
		Find(func(each ComponentTimeTableStruct) bool {
			return each.CalDate == today && each.IsWorkingDay == 1
		}).Result()
	if result != nil {
		return true
	}

	return false
}

func getLocation(authToken string) ComponentPosition {
	
}
