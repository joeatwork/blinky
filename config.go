package main

type configuration struct {
	device string
	deviceAddress string
	servicePort string
}

var (
	config = configuration{
		device: "1",
		deviceAddress: "0x09",
		servicePort: ":8080",
	}
)
