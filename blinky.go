package main

import (
	"fmt"
	"os"
	"time"
)

func serveBlinkM(device Device) (colors chan uint32, kill chan bool) {
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
				setColor(device, color)
			case die = <- kill:
				fmt.Printf("Blinkm qutting\n")
			}
		}
	}()
	return
}

func main() {
	fmt.Printf("TODO:\n")
	fmt.Printf("   SET UP POLLING GOROUTINE WITH RATE LIMITING\n")
	fmt.Printf("   FIGURE OUT WHICH I2C PART IS BUSTED AND REPLACE IT\n")

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s config_path\n", os.Args[0])
		os.Exit(1)
	}

	config, cerr := readConfig(os.Args[1])
	if cerr != nil {
		fmt.Printf("Can't read config file %s\n", os.Args[1])
		os.Exit(1)
	}

	var colorBlinkM chan<- uint32
	var killBlinkM chan<- bool
	colorBlinkM, killBlinkM = serveBlinkM(config.Device)
	defer func() {
		killBlinkM <- true
	}()
	go Run(config.ServicePort, colorBlinkM)

	redClient := InitClient(config.RedQuery.Token, config.RedQuery.Secret)
	greenClient := InitClient(config.GreenQuery.Token, config.GreenQuery.Secret)
	blueClient := InitClient(config.BlueQuery.Token, config.BlueQuery.Secret)

	qvalue, err := poll(config.RedQuery.Event, redClient, config.RedQuery.Where)
	if err != nil {
		panic("Bad response (or parse error) from Mixpanel")
	}
	fmt.Printf("queries value: %v\n", qvalue)
	mvalue, err := poll(config.GreenQuery.Event, greenClient, config.GreenQuery.Where)
	if err != nil {
		panic("Bad response (or parse error) from Mixpanel")
	}
	fmt.Printf("mobile value: %v\n", mvalue)
	nmvalue, err := poll(config.BlueQuery.Event, blueClient, config.BlueQuery.Where)
	if err != nil {
		panic("Bad response (or parse error) from Mixpanel")
	}
	fmt.Printf("non-mobile value: %v\n", nmvalue)

	fmt.Printf("%v, %v, %v\n",
		qvalue, mvalue, nmvalue)
	fmt.Printf("%x, %x, %x\n",
		(uint32(qvalue * 255) & 0xFF) << 16,
		(uint32(mvalue * 255) & 0xFF) << 8,
		uint32(nmvalue * 255) & 0xFF)
	qColor := (uint32(qvalue * 255) & 0xFF) << 16
	mColor := (uint32(mvalue * 255) & 0xFF) << 8
	nmColor := uint32(nmvalue * 255) & 0xFF

	color := qColor | mColor | nmColor

	fmt.Printf("COLOR: %x\n", color)
	// select {}
}
