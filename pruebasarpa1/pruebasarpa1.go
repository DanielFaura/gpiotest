package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	driver "gitlab.com/gomidi/rtmididrv"
	// driver "gitlab.com/gomidi/portmididrv"
)

var (
	transporte    byte     = 48
	modo          int      = 0
	cambioOctavas byte     = 12
	modos         [][]byte = [][]byte{
		{0, 2, 4, 5, 7, 9, 11}, // Mayor
		{0, 2, 4, 5, 7, 9, 10}, // Mayor con séptima
		{0, 2, 3, 5, 7, 8, 10}, // Menor natural
		{0, 2, 3, 5, 7, 8, 11}, // Menor armónica
		{0, 2, 3, 5, 7, 9, 11}, // Menor melódica
		{0, 2, 3, 6, 7, 8, 11}, // Escala árabe turbia
		{0, 2, 3, 5, 7, 9, 10}, // Dórica
	}
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func grado(cuerda byte) (grado byte) {
	// c = número de cuerda
	// g = grado de la escala
	grado =
		cuerda % 7
	return
}

func octava(cuerda byte) (octava byte) {
	// c = número de cuerda
	// g = grado de la escala
	octava = cuerda / 7
	return
}

// func nota(grado byte) (nota byte) {
// 	if modo == 0 {
// 		switch grado {
// 		case 0:
// 			nota = 0
// 		case 1:
// 			nota = 2
// 		case 2:
// 			nota = 4
// 		case 3:
// 			nota = 5
// 		case 4:
// 			nota = 7
// 		case 5:
// 			nota = 9
// 		case 6:
// 			nota = 11
// 		}
// 	}
// 	return
// }

func lector(pos *reader.Position, msg midi.Message) {
	fmt.Printf("got %s\n", msg)
	m := msg.Raw()
	fmt.Println("RECIBIDO", m)
	if m[0] == 0b10010000 && m[1] >= 48 && m[1] <= 59 && m[2] > 0 {
		// Cambio de tonalidad
		transporte = m[1] - 0
	}
	if m[0] == 0b10010000 && m[1] == 72 && m[2] > 0 {
		// Cambio de octava abajo
		cambioOctavas = cambioOctavas - 12
	}
	if m[0] == 0b10010000 && m[1] == 74 && m[2] > 0 {
		// Cambio de octava arriba
		cambioOctavas = cambioOctavas + 12
	}
	if m[0] == 0b10010000 && m[1] >= 60 && m[1] < 60+byte(len(modos)) && m[2] > 0 {
		// Cambio de modo
		modo = int(m[1]) - 60
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

	for _, v := range ins {
		fmt.Println("OUT", v)
	}

	in, out := ins[2], outs[1]

	check(in.Open())
	check(out.Open())

	defer in.Close()
	defer out.Close()

	noteOn := []byte{0b10010001, 0, 0}
	// noteOff := []byte{0b10000001, 0, 127}
	allSoundOff := []byte{0b10110001, 120, 0}

	out.Write(allSoundOff)
	time.Sleep(1 * time.Second)

	err = rpio.Open()
	check(err)

	rd := reader.New(
		reader.NoLogger(),
		reader.Each(lector),
	)

	err = rd.ListenTo(in)
	check(err)

	pines := make([]rpio.Pin, 22)

	for i := range pines {
		pines[i] = rpio.Pin(i + 4)
		pines[i].Input()
		pines[i].PullDown()
	}
	anterior := make([]rpio.State, 22)
	for {
		for i := range pines {
			if pines[i].Read() == rpio.Low && anterior[i] == rpio.High {
				// noteOn[1] = nota(grado(byte(i))) + 12*octava(byte(i)) + transporte + cambioOctavas
				noteOn[1] = modos[modo][grado(byte(i))] + 12*octava(byte(i)) + transporte + cambioOctavas
				noteOn[2] = 127
				fmt.Println("Cuerda:", i, ", grado:", grado(byte(i)), ", nota:", noteOn[1], ", transporte:", transporte, ", octava:", cambioOctavas)
				out.Write(noteOn)
				time.Sleep(10000 * time.Microsecond)
				anterior[i] = rpio.Low
			}
			if pines[i].Read() == rpio.High && anterior[i] == rpio.Low {
				// noteOn[1] = nota(grado(byte(i))) + 12*octava(byte(i)) + transporte + cambioOctavas
				noteOn[1] = modos[modo][grado(byte(i))] + 12*octava(byte(i)) + transporte + cambioOctavas
				noteOn[2] = 0
				out.Write(noteOn)
				time.Sleep(10000 * time.Microsecond)
				out.Write(noteOn)
				time.Sleep(10000 * time.Microsecond)
				anterior[i] = rpio.High
			}
		}
	}

}

//nota(grado(byte(i))) + 60
