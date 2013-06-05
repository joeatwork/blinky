package main

import (
	"fmt"
	"os"
	"time"
)

func serveBlinkM() (colors chan uint32, kill chan bool) {
	colors = make(chan uint32)
	kill = make(chan bool)
	go func () {
		die := false
		color := uint32(0)
		for ! die {
			time.Sleep(time.Second / 2)
			select {
			case color = <- colors:
				fmt.Printf("Color to %x\n", color)
				setColor(color)
			case die = <- kill:
				fmt.Printf("Qutting\n")
			}
		}
	}()
	return
}

func main() {
	if len(os.Args) != 4 && len(os.Args) != 1 {
		fmt.Printf("Usage: %s [R G B]\n", os.Args[0])
		os.Exit(1)
	}

	var colorBlinkM chan<- uint32
	var killBlinkM chan<- bool
	colorBlinkM, killBlinkM = serveBlinkM()
	defer func() {
		killBlinkM <- true
	}()

	if len(os.Args) == 4 {
		color, err := parseColor(os.Args[1], os.Args[2], os.Args[3])
		if err != nil {
			fmt.Printf("Can't parse color: %s\n", err.Error())
			os.Exit(1)
		}

		colorBlinkM <- color
	}
	go Run(colorBlinkM)

	select {}
}
