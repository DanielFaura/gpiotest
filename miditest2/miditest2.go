package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
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
		fmt.Println("OUT", v)
	}

	in, out := ins[1], outs[2]

	check(in.Open())
	check(out.Open())

	defer in.Close()
	defer out.Close()

	noteOff := []byte{0b10000001, 60, 127}
	noteOn := []byte{0b10010001, 60, 127}

	allSoundOff := []byte{0b10110001, 120, 0}

	out.Write(allSoundOff)
	time.Sleep(1 * time.Second)

	rd := reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Printf("got %s\n", msg)
			m := msg.Raw()
			fmt.Println("RECIBIDO", m)
			if m[0] == 0b10010000 && m[2] > 0 {
				noteOn[1] = m[1]
				fmt.Println("GUAPO", noteOn)
				out.Write(noteOn)
			}
			if m[0] == 0b10010000 && m[2] == 0 {
				noteOff[1] = m[1]
				fmt.Println("FEO", noteOff)
				out.Write(noteOff)
			}
		}),
	)

	err = rd.ListenTo(in)
	check(err)

	time.Sleep(1 * time.Hour)

	os.Exit(0)
}
