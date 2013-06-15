package main

import (
	"fmt"
	"os"
	"time"
)

func serveBlinkM(device Device) (colors chan uint32, kill chan bool) { // TODO move
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
	fmt.Printf("   REPORT COLORS TO WEB SERVER TO PREVIEW VARIATIONS OVER TIME\n")
	fmt.Printf("   GET A NEW RASPBERRY PI BOARD WITH WORKING I2C\n")

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
	go RunWebService(config.ServicePort, colorBlinkM)

	redClient := InitClient(config.RedQuery.Token, config.RedQuery.Secret)
	greenClient := InitClient(config.GreenQuery.Token, config.GreenQuery.Secret)
	blueClient := InitClient(config.BlueQuery.Token, config.BlueQuery.Secret)
	pollTime := config.PollingRateSeconds * time.Second
	if pollTime < 30 {
		fmt.Printf("Polling time must be 30 seconds or more.")
		os.Exit(1)
	}
	redVal, redKill := RunPollingService(
		pollTime,
		config.RedQuery.Event,
		config.RedQuery.Where,
		redClient,
	)
	defer func() {
		redKill <- true
	}()
	greenVal, greenKill := RunPollingService(
		pollTime,
		config.GreenQuery.Event,
		config.GreenQuery.Where,
		greenClient,
	)
	defer func() {
		greenKill <- true
	}()
	blueVal, blueKill := RunPollingService(
		pollTime,
		config.BlueQuery.Event,
		config.BlueQuery.Where,
		blueClient,
	)
	defer func() {
		blueKill <- true
	}()
	rok := true
	gok := true
	bok := true
	var r, g, b float64
	for rok && gok && bok {
		fmt.Printf("Waiting for pollers...\n")
		select {
		case r, rok = <- blueVal:
		case g, gok = <- greenVal:
		case b, bok = <- redVal:
		}

		if rok && gok && bok {
			fmt.Printf(" %v, %v, %v\n", r, g, b)
			rColor := (uint32(r * 255) & 0xFF) << 16
			gColor := (uint32(g * 255) & 0xFF) << 8
			bColor := uint32(b * 255) & 0xFF

			color := rColor | gColor | bColor

			fmt.Printf("COLOR: %x\n", color)
		} else {
			fmt.Printf(" Poller died! red %v green %v blue %v\n",
				rok, gok, bok)
		}
	}
}
