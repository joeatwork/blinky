package main

import (
	"fmt"
	"net/http"
)

const (
	colorForm = `<html>
<body>
<form method="POST">
    <div>r: <input type="text" name="r"></div>
    <div>g: <input type="text" name="g"></div>
    <div>b: <input type="text" name="b"></div>
    <div><input type="submit"></div>
</form>
</body>
</html>
`
)

type blinkMHandler struct {
	blinkM chan<- uint32
}

func (handler *blinkMHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rs := r.FormValue("r")
	gs := r.FormValue("g")
	bs := r.FormValue("b")

	if rs == "" || gs == "" || bs == "" {
		h := w.Header()
		h.Set("Content-Type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, colorForm)
		return
	}

	color, err := parseColor(rs, gs, bs)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Can't understand r,g,b as integers/color")
		return
	}

	fmt.Fprintf(w, "1")
	handler.blinkM <- color
	return
}

func RunWebService(servicePort string, colorBlinkM chan<- uint32) {
	handler := &blinkMHandler{blinkM: colorBlinkM}
	fmt.Printf("Starting service at %s\n", servicePort)
	http.Handle("/", handler)
	err := http.ListenAndServe(servicePort, nil)
	if err != nil {
		fmt.Printf("Couldn't serve: %s\n", err.Error())
	}
}
