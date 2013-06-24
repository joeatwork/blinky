package main

import (
	"fmt"
	"os/exec"
)

func runBlinkmScript(device Device) {
	// Stop any currently running script
	stop := exec.Command("i2cset", "-y", device.Device, device.DeviceAddress, "0x6f")
	stoperr := stop.Run()
	if stoperr != nil {
		fmt.Printf("OH NO! %s\n", stoperr.Error())		
	}

	// Run the SOS script 3 times from start
	sosScript := exec.Command("i2cset", "-y", device.Device, device.DeviceAddress, "0x70", "0x12", "0x03", "0x00", "i")	
	soserr := sosScript.Run()
	if soserr != nil {
		fmt.Printf("OH NO! %s\n", soserr.Error())
	}
}

// See http://thingm.com/fileadmin/thingm/downloads/BlinkM_datasheet.pdf
func setColor(device Device, color uint32) {
	r := fmt.Sprintf("0x%X", (color >> 16) & 0xFF)
	g := fmt.Sprintf("0x%X", (color >> 8) & 0xFF)
	b := fmt.Sprintf("0x%X", color & 0xFF)
	fmt.Printf("R %s G %s B %s\n", r, g, b)

	// Stop any currently running script
	stop := exec.Command("i2cset", "-y", device.Device, device.DeviceAddress, "0x6f")
	stoperr := stop.Run()
	if stoperr != nil {
		fmt.Printf("OH NO! %s\n", stoperr.Error())		
	}

	// Set to slowest fade speed
	fadeSpeed := exec.Command("i2cset", "-y", device.Device, device.DeviceAddress, "0x66", "0x01", "i")
	fadeerr := fadeSpeed.Run()
	if fadeerr != nil {
		fmt.Printf("OH NO! %s\n", fadeerr.Error())
	}

	// Fade to RGB color
	cmd := exec.Command("i2cset", "-y", device.Device, device.DeviceAddress, "0x63", r, g, b, "i")
	err := cmd.Run()
	if err != nil {
	 	fmt.Printf("OH NO! %s\n", err.Error())
	}
}
