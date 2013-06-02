package main

import (
	"fmt"
	"os/exec"
)

func setColor(color uint32) {
	r := fmt.Sprintf("0x%X", (color >> 16) & 0xFF)
	g := fmt.Sprintf("0x%X", (color >> 8) & 0xFF)
	b := fmt.Sprintf("0x%X", color & 0xFF)
	fmt.Printf("R %s G %s B %s\n", r, g, b)

	stop := exec.Command("i2cset", "-y", config.device, config.deviceAddress, "0x6f")
	stoperr := stop.Run()
	if stoperr != nil {
		fmt.Printf("OH NO! %s\n", stoperr.Error())		
	}

	cmd := exec.Command("i2cset", "-y", config.device, config.deviceAddress, "0x63", r, g, b, "i")
	err := cmd.Run()
	if err != nil {
	 	fmt.Printf("OH NO! %s\n", err.Error())
	}
}
