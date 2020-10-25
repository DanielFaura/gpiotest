package main

import (
	"fmt"
	// driver "gitlab.com/gomidi/portmididrv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	drv, err := drv.New()
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

}
