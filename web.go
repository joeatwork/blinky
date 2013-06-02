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

func handler(w http.ResponseWriter, r *http.Request) {
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

	setColor(color)
	fmt.Fprintf(w, "Set 0x%x", color)
	return
}

func Run() {
	fmt.Printf("Planning to run a server\n")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Couldn't serve: %s\n", err.Error())
	}
}