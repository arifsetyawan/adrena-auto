package main

import (
	"arisetyawan/adrena-auto/hatcher"
	"fmt"
	"os"
)

var authToken = ""

func main() {
	fmt.Println("Adrena Automation")
	authToken, _ = hatcher.Auth()
	isWorkingDay := hatcher.IsTodayWorkingOrNot(authToken)
	if isWorkingDay == false {
		fmt.Println("NOT WORKING DAY")
		os.Exit(1)
	}

}
