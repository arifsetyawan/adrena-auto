package hatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/novalagung/gubrak/v2"
)

var baseURL = "https://adrena.tech"

type BridgeInterface interface {
	Auth()
	IsWorking() bool
	GetLocation()
	DoCheck(activityType string) bool
	HttpWrap(params HttpWrapParams) ([]byte, error)
}

type Bridge struct {
	BridgeInterface
	AuthToken string
	Position  ComponentPosition
}

type ComponentPosition struct {
	LocationName string `json:"locationName"`
	ProvinceId   int    `json:"provinceId"`
}

func NewBridge() BridgeInterface {
	return &Bridge{}
}

func (b *Bridge) Auth() {

	var data ResponseAuth

	if os.Getenv("AUTH_USERNAME") == "" || os.Getenv("AUTH_PASSWORD") == "" {
		fmt.Println("NO CREDENTIAL")
		os.Exit(1)
	}

	values := map[string]string{"username": os.Getenv("AUTH_USERNAME"), "password": os.Getenv("AUTH_PASSWORD")}
	bodyJSON, _ := json.Marshal(values)
	var payload = bytes.NewBuffer(bodyJSON)

	body, errHttp := b.HttpWrap(HttpWrapParams{method: "POST", url: baseURL + "/auth/api/mobile/session", payload: payload})
	if errHttp != nil {
		fmt.Println("ERROR HTTP")
		os.Exit(1)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	fmt.Println(string(body))
	fmt.Println(data)
	fmt.Println("TOKEN", data.Token.Access)

	b.AuthToken = data.Token.Access
}

func (b *Bridge) IsWorking() bool {
	var data ResponseTimeTable

	// print location and local time
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println(err)
	}

	currentTime := time.Now().In(location)
	today := currentTime.Format("2006-01-02")

	fmt.Println("LOCATION: " + location.String())
	fmt.Println("today: " + today)

	fmt.Println(baseURL + "/ess/api/timetable/published/week?date=" + today)

	body, errHttp := b.HttpWrap(HttpWrapParams{method: "GET", url: baseURL + "/ess/api/timetable/published/week?date=" + today, payload: nil, useAuth: true})
	if errHttp != nil {
		fmt.Println("ERROR HTTP")
		os.Exit(1)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	fmt.Println(data)

	result := gubrak.From(data.TimeTable).
		Find(func(each ComponentTimeTableStruct) bool {
			return each.CalDate == today && each.IsWorkingDay == 1
		}).Result()

	if result != nil {
		return true
	}

	return false
}

func (b Bridge) GetLocation() {
	var data ResponseLocation

	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")
	lat := os.Getenv("LAT")
	long := os.Getenv("LONG")

	body, errHttp := b.HttpWrap(HttpWrapParams{method: "GET", url: baseURL + "/ess/api/attendance/position/v2?date=" + today + "&lat=" + lat + "&long=" + long, payload: nil, useAuth: true})
	if errHttp != nil {
		fmt.Println("ERROR HTTP")
		os.Exit(1)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	b.Position = data.Position
}

func (b *Bridge) DoCheck(activityType string) bool {

	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")
	lat, _ := strconv.ParseFloat(os.Getenv("LAT"), 128)
	long, _ := strconv.ParseFloat(os.Getenv("LONG"), 128)

	b.Auth()

	if !b.IsWorking() {
		fmt.Println("NOT WORKING DAY")
		os.Exit(1)
	}

	fmt.Println("WORKING DAY")

	b.GetLocation()

	values := map[string]interface{}{
		"activityType":   activityType,
		"clockInMethod":  "GEOLOC",
		"deviceId":       -1,
		"latitude":       lat,
		"longitude":      long,
		"locationName":   b.Position.LocationName,
		"provinceId":     b.Position.ProvinceId,
		"selectedDate":   today,
		"workLocationId": -1,
	}
	bodyJSON, _ := json.Marshal(values)
	var payload = bytes.NewBuffer(bodyJSON)

	body, errHttp := b.HttpWrap(HttpWrapParams{method: "POST", url: baseURL + "/ess/api/attendance", payload: payload, useAuth: true})
	if errHttp != nil {
		fmt.Println("ERROR HTTP")
		os.Exit(1)
	}

	if body != nil {
		fmt.Println(string(body))
		return true
	}

	return false

}

type HttpWrapParams struct {
	method  string
	url     string
	payload io.Reader
	useAuth bool
}

func (b *Bridge) HttpWrap(params HttpWrapParams) ([]byte, error) {
	var client = &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	req, err := http.NewRequest(params.method, params.url, params.payload)
	if err != nil {
		fmt.Println("ERROR AUTH: CANNOT CONNECT TO SERVER")
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "Adrena%20HR/1 CFNetwork/1325.0.1 Darwin/21.1.0")
	if params.method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}

	if params.useAuth {
		fmt.Println("AUTH_TOKEN", b.AuthToken)
		req.Header.Set("Authorization", "Bearer "+b.AuthToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(body))

	return body, nil
}

type ResponseAuth struct {
	Token         ResponseAuthToken `json:"token"`
	ResultSuccess int               `json:"resultSuccess"`
}
type ResponseAuthToken struct {
	Access  string `json:"access"`
	Id      string `json:"id"`
	Refresh string `json:"refresh"`
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
