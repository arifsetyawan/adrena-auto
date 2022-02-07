package main

import (
	"arisetyawan/adrena-auto/hatcher"
	"fmt"
)

func main() {
	fmt.Println("Adrena Automation Check IN")
	hatcherBridge := hatcher.NewBridge()
	hatcherBridge.DoCheck("CHECKIN")
}
