package main

import (
	"fmt"
	"time"
)

func main() {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println(err)
	}

	currentTime := time.Now().In(location)
	today := currentTime.Format("2006-01-02")
	fmt.Println(today)
}
