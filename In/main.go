package main

import (
	"arisetyawan/adrena-auto/hatcher"
	"fmt"
	"os"
)

var authToken = ""

func main() {
	fmt.Println("Adrena Automation Check IN")
	authToken, _ = hatcher.Auth()
	isWorkingDay := hatcher.IsTodayWorkingOrNot(authToken)
	if isWorkingDay == false {
		fmt.Println("NOT WORKING DAY")
		os.Exit(1)
	}
	position := hatcher.GetLocation(authToken)
	fmt.Println(position)
	hatcher.Check(authToken, "CHECKIN", position)
	fmt.Println("AUTO CHECKIN DONE ")
}
