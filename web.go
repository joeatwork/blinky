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

    .sample_bar {
        height: 5;
        background-color: black;
    }
    </style>
</head>
<body>
    <div style="width: 300px;">
        {{ range .Moments }}
        <div style="background-color: #{{ printf "%06x" .Color }};"
             class="sample_demo"
             >{{ .Time }}</div>
        <div style="width: {{ .RPercent }}%; background-color: red;"
             class="sample_bar"></div>
        <div style="width: {{ .GPercent }}%; background-color: green;"
             class="sample_bar"></div>
        <div style="width: {{ .BPercent }}%; background-color: blue;"
             class="sample_bar"></div>
        {{ end }}
    </div><!-- container -->
</body>
</html>
`
)

type colorMoment struct {
	Color    uint32
	RPercent uint32
	GPercent uint32
	BPercent uint32
	Time     time.Time
}

type service struct {
	Moments  []colorMoment
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

func colorPercents(color uint32) (r uint32, g uint32, b uint32) {
	r = (((color >> 16) & 0xFF) * 100) / 256
	g = (((color >> 8) & 0xFF) * 100) / 256
	b = ((color & 0xFF) * 100) / 256
	return
}

func RunWebService(servicePort string, colors <-chan uint32) {
	samples := make([]colorMoment, 50)
	t, err := template.New("Moments").Parse(colorForm)
	if err != nil {
		fmt.Printf("Can't interpret template: %s\n", err.Error())
		os.Exit(1)
	}
	var mutex sync.RWMutex
	service := &service{samples, t, mutex}
	fmt.Printf("Starting service at %s\n", servicePort)
	http.Handle("/", service)

	go func() {
		open := true
		var color uint32 = 0
		var index = 0
		for open {
			color, open = <-colors
			if open {
				rPercent, gPercent, bPercent := colorPercents(color)
				now := time.Now()
				service.Lock()
				service.Moments[index] = colorMoment{
					color,
					rPercent,
					gPercent,
					bPercent,
					now,
				}
				service.Unlock()
				index = index + 1
				if index >= len(service.Moments) {
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
