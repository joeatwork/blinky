package main

type configuration struct {
	device string
	deviceAddress string
}

var (
	config = configuration{
		device: "1",
		deviceAddress: "0x09",
	}
)