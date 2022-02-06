package main

import (
	"arisetyawan/adrena-auto/hatcher"
	"fmt"
	"os"
)

var authToken = ""

func main() {
	fmt.Println("Adrena Automation Check Out")
	authToken, _ = hatcher.Auth()
	isWorkingDay := hatcher.IsTodayWorkingOrNot(authToken)
	if isWorkingDay == false {
		fmt.Println("NOT WORKING DAY")
		os.Exit(1)
	}
	position := hatcher.GetLocation(authToken)
	fmt.Println(position)
	hatcher.Check(authToken, "CHECKOUT", position)
	fmt.Println("AUTO CHECKOUT DONE ")
}
