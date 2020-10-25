package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {
	err := rpio.Open()
	if err != nil {
		log.Panic("Error!")
	}
	placa := rpio.Pin(3)
	placa.Input()
	led := rpio.Pin(2)
	led.Output()
	led.Low()
	for {
		fmt.Println(placa.Read())
		time.Sleep(200 * time.Millisecond)
		if placa.Read() == 0 {
			led.High()
		}
		if placa.Read() == 1 {
			led.Low()
		}
	}
}
