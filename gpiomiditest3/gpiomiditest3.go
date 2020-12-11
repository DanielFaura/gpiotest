package main

import (
	"fmt"
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

func grado(cuerda byte) (grado byte) {
	// c = número de cuerda
	// g = grado de la escala
	grado = (cuerda - 4) % 7
	return
}

func nota(grado byte) (nota byte) {
	switch grado {
	case 0:
		nota = 0
	case 1:
		nota = 2
	case 2:
		nota = 4
	case 3:
		nota = 5
	case 4:
		nota = 7
	case 5:
		nota = 9
	case 6:
		nota = 12
	}
	return
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

	for _, v := range ins {
		fmt.Println("OUT", v)
	}

	in, out := ins[1], outs[1]

	check(in.Open())
	check(out.Open())

	defer in.Close()
	defer out.Close()

	noteOn := []byte{0b10010001, 0, 127}
	noteOff := []byte{0b10000001, 0, 127}
	allSoundOff := []byte{0b10110001, 120, 0}

	out.Write(allSoundOff)
	time.Sleep(1 * time.Second)

	err = rpio.Open()
	check(err)

	pines := make([]rpio.Pin, 25)

	for i := range pines {
		pines[i] = rpio.Pin(i + 1)
		pines[i].Input()
		pines[i].PullDown()
	}
	anterior := make([]rpio.State, 25)
	for {
		for i := range pines {
			if pines[i].Read() == rpio.Low && anterior[i] == rpio.High {
				fmt.Println("XD SUENA")
				noteOn[1] = nota(grado(byte(i))) + 60
				out.Write(noteOn)
				time.Sleep(2000 * time.Microsecond)
				anterior[i] = rpio.Low
			}
			if pines[i].Read() == rpio.High && anterior[i] == rpio.Low {
				fmt.Println("XD NO SUENA")
				noteOff[1] = nota(grado(byte(i))) + 60
				out.Write(noteOff)
				time.Sleep(2000 * time.Second)
				anterior[i] = rpio.High
			}
		}
	}

}
