package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
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
             >{{ .When }}</div>
        <div style="width: {{ percent .RedScaled }}%; background-color: red;"
             class="sample_bar"></div>
        <div style="width: {{ percent .GreenScaled }}%; background-color: green;"
             class="sample_bar"></div>
        <div style="width: {{ percent .BlueScaled }}%; background-color: blue;"
             class="sample_bar"></div>
        {{ end }}
    </div><!-- container -->
</body>
</html>
`
)

type service struct {
	Template  *template.Template
	LogSource func() []LogSample
}

type report struct {
	Moments []LogSample
}

func percent(f float64) string {
	return fmt.Sprintf("%f", f*100)
}

func (service *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/html")
	w.WriteHeader(200)

	rpt := report{service.LogSource()}
	err := service.Template.Execute(w, rpt)
	if err != nil {
		fmt.Printf("Can't execute web service template: %s\n", err.Error())
	}
}

func RunWebService(servicePort string, logger *Logger) {
	t := template.New("Moments")
	t.Funcs(template.FuncMap{"percent": percent})
	t, err := t.Parse(colorForm)
	if err != nil {
		fmt.Printf("Can't interpret template: %s\n", err.Error())
		os.Exit(1)
	}

	recentService := &service{t, logger.recent}
	ancientService := &service{t, logger.ancient}
	fmt.Printf("Starting service at %s\n", servicePort)
	http.Handle("/", recentService)
	http.Handle("/history", ancientService)

	go func() {
		err := http.ListenAndServe(servicePort, nil)
		if err != nil {
			fmt.Printf("Couldn't serve: %s\n", err.Error())
		}
	}()
}
