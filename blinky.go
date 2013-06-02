package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s R G B\n", os.Args[0])
		os.Exit(1)
	}
	color, err := parseColor(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		fmt.Printf("Can't parse color: %s\n", err.Error())
		os.Exit(1)
	}

	setColor(color);
	Run()
}
