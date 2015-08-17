package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bugsnag/bugsnag-go"
)

var (
	port      int
	redirects []Redirect
)

func init() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       os.Getenv("BUGSNAG_API_KEY"),
		ReleaseStage: os.Getenv("BUGSNAG_RELEASE_STAGE"),
	})

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
		redirect, err := FindRedirectForTwilio(redirects, number)
		if err != nil {
			responseCode = http.StatusNotFound
		} else {
			responseCode = http.StatusOK
			w.Write(GenerateResponseXMLFor(redirect))
		}
	}

	w.WriteHeader(responseCode)

	fmt.Printf("%s %s %d\n", r.Method, r.URL.Path, responseCode)
}

func LookupHandler(w http.ResponseWriter, r *http.Request) {
	var responseCode int

	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/lookup/")

	if id == "" {
		responseCode = http.StatusBadRequest
	} else {
		redirectTo, err := FindTwilioForID(redirects, id)
		if err != nil {
			responseCode = http.StatusNotFound
		} else {
			responseCode = http.StatusOK
			w.Write(GenerateResponseJSONFor(redirectTo))
		}
	}

	w.WriteHeader(responseCode)

	fmt.Printf("%s %s %d\n", r.Method, r.URL.Path, responseCode)
}

func main() {
	fmt.Printf("> Starting on http://0.0.0.0:%d\n", port)

	http.HandleFunc("/twilio", TwilioHandler)
	http.HandleFunc("/lookup/", LookupHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), bugsnag.Handler(nil))
	if err != nil {
		panic("Error starting!")
	}
}
