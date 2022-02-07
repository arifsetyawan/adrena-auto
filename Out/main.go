package main

import (
	"arisetyawan/adrena-auto/hatcher"
	"fmt"
)

func main() {
	fmt.Println("Adrena Automation Check Out")
	hatcherBridge := hatcher.NewBridge()
	hatcherBridge.DoCheck("CHECKOUT")
}
