package main

import (
	"log"

	"github.com/stianeikeland/go-rpio"

	"time"
)

func main() {
	err := rpio.Open()
	if err != nil {
		log.Panic("Error!")
	}
	led := rpio.Pin(2)
	led.Output()
	led.Low()
	pulsador := rpio.Pin(4)
	pulsador.Input()
	pulsador.PullDown()
	anterior := pulsador.Read()
	for {
		actual := pulsador.Read()
		if anterior == 0 && actual == 1 {
			led.Toggle()
			time.Sleep(200 * time.Millisecond)
		}
		anterior = actual
	}
}
