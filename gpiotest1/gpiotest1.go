package main

import (
	"log"

	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {
	err := rpio.Open()
	led := rpio.Pin(2)
	led.Output()
	if err != nil {
		log.Panic("Errorcillo!")
	}
	led.Low()
	for {
		time.Sleep(time.Second / 16)
		led.High()
		time.Sleep(time.Second / 16)
		led.Low()
	}

}
