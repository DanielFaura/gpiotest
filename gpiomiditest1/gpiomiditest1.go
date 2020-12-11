package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
	driver "gitlab.com/gomidi/rtmididrv"
	// driver "gitlab.com/gomidi/portmididrv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	drv, err := driver.New()
	check(err)

	defer drv.Close()

	ins, err := drv.Ins()
	check(err)

	outs, err := drv.Outs()
	check(err)

	for _, v := range ins {
		fmt.Println("IN", v)
	}

	for _, v := range outs {
		fmt.Println("OUTS", v)
	}

	in, out := ins[0], outs[1]

	check(in.Open())
	check(out.Open())

	defer in.Close()
	defer out.Close()

	// noteOn := []byte{0b10010001, 0, 127}
	// noteOff := []byte{0b10000001, 0, 127}
	allSoundOff := []byte{0b10110001, 120, 0}

	out.Write(allSoundOff)
	time.Sleep(1 * time.Second)

	err = rpio.Open()
	check(err)

	pines := make([]rpio.Pin, 22)
	// pines[0] = rpio.Pin(1)
	// pines[0].Input()
	// pines[0].PullDown()
	for i := range pines {
		pines[i] = rpio.Pin(i + 1)
		pines[i].Input()
		pines[i].PullDown()
	}

	// anterior := 0

	for {
		fmt.Println("Leyendo pin, el valor es", pines[0].Read())
		time.Sleep(1 * time.Second)
	}

	// for {
	// 	if placa.Read() == 0 && anterior == 0 {
	// 		noteOn[1] = 60
	// 		out.Write(noteOn)
	// 		time.Sleep(2000 * time.Microsecond)
	// 		anterior = 1
	// 	}
	// 	if placa.Read() == 1 {
	// 		noteOff[1] = 60
	// 		out.Write(noteOff)
	// 		time.Sleep(2000 * time.Microsecond)
	// 		anterior = 0
	// 	}
	// }

	time.Sleep(1 * time.Hour)

	os.Exit(0)

}
