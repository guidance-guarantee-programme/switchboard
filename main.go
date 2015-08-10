package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
)

var (
	port      int
	redirects []Redirect
)

func init() {
	flag.IntVar(&port, "p", 8080, "Port to listen on")
	flag.Parse()

	redirects = LoadRedirectsFromYAML("redirects.yaml")
}

func TwilioHandler(w http.ResponseWriter, r *http.Request) {
	var responseCode int

	number := r.FormValue("To")

	if number == "" {
		responseCode = http.StatusBadRequest
	} else {
		redirectTo, err := FindCabForTwilio(redirects, number)
		if err != nil {
			responseCode = http.StatusNotFound
		} else {
			responseCode = http.StatusOK
			w.Write(GenerateResponseXMLFor(redirectTo))
		}
	}

	w.WriteHeader(responseCode)

	fmt.Printf("%s %s %d\n", r.Method, r.URL.Path, responseCode)
}

func LookupHandler(w http.ResponseWriter, r *http.Request) {
	var responseCode int

	id := strings.TrimPrefix(r.URL.Path, "/lookup/")

	if id == "" {
		responseCode = http.StatusBadRequest
	} else {
		redirectTo, err := FindTwilioForID(redirects, id)
		if err != nil {
			responseCode = http.StatusNotFound
		} else {
			responseCode = http.StatusOK
			w.Write([]byte(redirectTo))
		}
	}

	w.WriteHeader(responseCode)

	fmt.Printf("%s %s %d\n", r.Method, r.URL.Path, responseCode)
}

func main() {
	fmt.Printf("> Starting on http://0.0.0.0:%d\n", port)

	http.HandleFunc("/twilio", TwilioHandler)
	http.HandleFunc("/lookup/", LookupHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic("Error starting!")
	}
}
