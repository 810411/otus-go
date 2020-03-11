package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const host = "ntp1.stratum1.ru"

func main() {
	currentTime := time.Now()
	exactTime, err := ntp.Time(host)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Printf("current time: %v\nexact time: %v\n", currentTime, exactTime)
}
