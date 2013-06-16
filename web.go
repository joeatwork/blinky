package main

import (	
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"
)

const (
	colorForm = `<html>
<head>
    <style>
    .sample_demo {
        width: 300px;
        height: 30px;
    }
    </style>
</head>
<body>
    <h1>YO DAWG HERE YOU ARE</h1>
    {{ range .Samples }}
    <div style="background-color: #{{ printf "%06x" .Color }};"
         class="sample_demo"
         >{{ .Time }}</div>
    {{ end }}
</body>
</html>
`
)

type sample struct {
	Color uint32
	Time time.Time
}

type service struct {
	Samples []sample
	Template *template.Template
	sync.RWMutex
}

func (service *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/html")
	w.WriteHeader(200)

	service.RLock()
	service.Template.Execute(w, service)
	service.RUnlock()
}

func RunWebService(servicePort string, colors <-chan uint32) {
	samples := make([]sample, 50)
	t, err := template.New("Samples").Parse(colorForm)
	if err != nil {
		fmt.Printf("Can't interpret template: %s\n", err.Error())
		os.Exit(1)
	}
	var mutex sync.RWMutex
	service := &service{ samples, t, mutex }
	fmt.Printf("Starting service at %s\n", servicePort)
	http.Handle("/", service)

	go func() {
		open := true
		var color uint32 = 0
		var index = 0
		for open {
			fmt.Printf("Web Server Waiting for Color\n")
			color, open = <-colors
			fmt.Printf("COLOR %60x OPEN %v\n", color, open)
			if open {
				now := time.Now()
				service.Lock()
				service.Samples[index] = sample{ color, now }
				service.Unlock()
				index = index + 1
				if index >= len(service.Samples) {
					index = 0
				}
			}
		}
	}()

	go func() {
		err := http.ListenAndServe(servicePort, nil)
		if err != nil {
			fmt.Printf("Couldn't serve: %s\n", err.Error())
		}
	}()
}
