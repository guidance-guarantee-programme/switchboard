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
	number := r.FormValue("To")

	if number == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redirect, err := FindRedirectForTwilio(redirects, number)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(GenerateResponseXMLFor(redirect))
	}
}

func LookupHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/lookup/")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redirectTo, err := FindTwilioForID(redirects, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(GenerateResponseJSONFor(redirectTo))
	}
}

func main() {
	fmt.Printf("> Starting on http://0.0.0.0:%d\n", port)

	http.HandleFunc("/twilio", TwilioHandler)
	http.HandleFunc("/lookup/", LookupHandler)

	handler := NewLoggingMiddleware(bugsnag.Handler(nil))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		panic("Error starting!")
	}
}
