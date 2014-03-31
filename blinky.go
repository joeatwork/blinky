package main

import (
	"fmt"
	"math"
	"math/rand"
	"flag"
	"os"
	"time"
)

func serveBlinkM(device Device) (colors chan uint32, kill chan bool) {
	colors = make(chan uint32)
	kill = make(chan bool)
	go func() {
		die := false
		color := uint32(0)
		for !die {
			time.Sleep(time.Second / 2)
			select {
			case color = <-colors:
				fmt.Printf("Color to %x\n", color)
				setColor(device, color)
			case die = <-kill:
				fmt.Printf("Blinkm qutting\n")
			}
		}
	}()
	return
}

func servePollColor(pollingRate time.Duration, event, where string, client *ReportClient) (value chan float64, kill chan bool) {
	// Closes the channel on error
	value = make(chan float64)
	kill = make(chan bool, 1)
	die := false
	go func() {
		time.Sleep(time.Second * time.Duration(rand.Intn(60)))
		for !die {
			samples, err := poll(event, where, client)
			if err != nil {
				fmt.Printf("ERROR IN POLLER")
				fmt.Printf("  %v, %v, %v", event, where, client)
				fmt.Println(err.Error())
				value <- 0
			} else {
				value <- colorForCurrentSample(samples)
			}

			time.Sleep(pollingRate)
			select {
			case die = <-kill:
			default:
			}
		}
		fmt.Printf("Poller exiting")
		close(value)
	}()

	return
}

func serveLog(logger *Logger) (ret chan LogSample) {
	ret = make(chan LogSample, 10)
	go func() {
		for log := range ret {
			logger.writeLog(log)
		}
	}()
	return
}

var debugMode bool

func init() {
	flag.BoolVar(&debugMode, "verbose", false, "Print lots of stuff to stdout")
}

func main() {
	flag.Parse()

	fmt.Printf("Debug mode == %t\n", debugMode)
	fmt.Printf("TODO:\n")
	fmt.Printf("   Smooth out artifacts near the top of the hour\n")

	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s config_path\n", os.Args[0])
		os.Exit(1)
	}

	config, cerr := readConfig(flag.Arg(0))
	if cerr != nil {
		fmt.Printf("Can't read config file %s\n", flag.Arg(0))
		os.Exit(1)
	}

	runBlinkmScript(config.Device)

	// BlinkM
	var colorBlinkM chan<- uint32
	var killBlinkM chan<- bool
	colorBlinkM, killBlinkM = serveBlinkM(config.Device)
	defer func() {
		killBlinkM <- true
	}()

	// Logger
	logger := initLog(debugMode)
	logChannel := serveLog(logger)

	RunWebService(config.ServicePort, logger)

	// Polling Mixpanel
	redClient := InitClient(config.RedQuery.Token, config.RedQuery.Secret)
	greenClient := InitClient(config.GreenQuery.Token, config.GreenQuery.Secret)
	blueClient := InitClient(config.BlueQuery.Token, config.BlueQuery.Secret)
	pollTime := config.PollingRateSeconds * time.Second
	if pollTime < 30 {
		fmt.Printf("Polling time must be 30 seconds or more.")
		os.Exit(1)
	}
	redVal, redKill := servePollColor(
		pollTime,
		config.RedQuery.Event,
		config.RedQuery.Where,
		redClient,
	)
	defer func() {
		redKill <- true
	}()
	greenVal, greenKill := servePollColor(
		pollTime,
		config.GreenQuery.Event,
		config.GreenQuery.Where,
		greenClient,
	)
	defer func() {
		greenKill <- true
	}()
	blueVal, blueKill := servePollColor(
		pollTime,
		config.BlueQuery.Event,
		config.BlueQuery.Where,
		blueClient,
	)
	defer func() {
		blueKill <- true
	}()

	// Messages - Poll to BlinkM
	rok := true
	gok := true
	bok := true
	var r, g, b float64
	for rok && gok && bok {
		select {
		case r, rok = <-blueVal:
		case g, gok = <-greenVal:
		case b, bok = <-redVal:
		}

		if rok && gok && bok {
			// We stretch the color scale a bit, for DRAMA
			rDrama := math.Pow(r, 1.2)
			gDrama := math.Pow(g, 1.2)
			bDrama := math.Pow(b, 1.2)

			rColor := (uint32(rDrama*255) & 0xFF) << 16
			gColor := (uint32(gDrama*255) & 0xFF) << 8
			bColor := uint32(bDrama*255) & 0xFF
			color := rColor | gColor | bColor
			colorBlinkM <- color

			logChannel <- LogSample{
				r, rDrama,
				g, gDrama,
				b, bDrama,
				color, time.Now(),
			}
		} else {
			fmt.Printf("Poller died! red %v green %v blue %v\n",
				rok, gok, bok)
		}
	}
}
