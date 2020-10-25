package main

import (
	"fmt"
	//"math/rand"
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

	in, out := ins[1], outs[1]

	check(in.Open())
	check(out.Open())

	defer in.Close()
	defer out.Close()

	//noteOff := []byte{0b10000001, 60, 127}
	//noteOn := []byte{0b10010001, 60, 127}
	allSoundOff := []byte{0b10110001, 120, 0}

	out.Write(allSoundOff)
	time.Sleep(1 * time.Second)

	rd := reader.New(
		reader.NoLogger(),
		// write every message to the out port
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Printf("got %s\n", msg)
		}),
	)

	err = rd.ListenTo(in)
	check(err)

	time.Sleep(1 * time.Hour)

	os.Exit(0)
}
