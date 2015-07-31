package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var responseXML = `
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Dial>789</Dial>
</Response>
`

func TestLoadRedirectsFromYAML(t *testing.T) {
	var (
		redirects []Redirect
		redirect  Redirect
	)

	redirects = LoadRedirectsFromYAML("fixtures/redirects.yaml")

	if len(redirects) != 1 {
		t.Error("Should return a slice of 1 Redirect")
	}

	redirect = redirects[0]

	if redirect.Id != "123" {
		t.Error("Id should be set")
	}

	if redirect.Twilio != "456" {
		t.Error("Twilio should be set")
	}

	if redirect.Cab != "789" {
		t.Error("Cab should be set")
	}
}

func TestFindCabForTwilio(t *testing.T) {
	redirects := LoadRedirectsFromYAML("fixtures/redirects.yaml")

	if _, err := FindCabForTwilio(redirects, "456"); err != nil {
		t.Error("Should not return an error for valid Twilio number")
	}

	if redirect, _ := FindCabForTwilio(redirects, "456"); redirect != "789" {
		t.Error("Should return CAB number for valid Twilio number")
	}

	if _, err := FindCabForTwilio(redirects, "999"); err == nil {
		t.Error("Should return an error for invalid Twilio number")
	}
}

func TestFindTwilioForID(t *testing.T) {
	redirects := LoadRedirectsFromYAML("fixtures/redirects.yaml")

	if _, err := FindTwilioForID(redirects, "123"); err != nil {
		t.Error("Should not return an error for valid id")
	}

	if redirect, _ := FindTwilioForID(redirects, "123"); redirect != "456" {
		t.Error("Should return Twilio number for valid id")
	}

	if _, err := FindTwilioForID(redirects, "000"); err == nil {
		t.Error("Should return an error for invalid id")
	}
}

func TestGenerateResponseXMLFor(t *testing.T) {
	expected := []byte(strings.TrimSpace(responseXML))
	result := GenerateResponseXMLFor("789")

	if !bytes.Equal(expected, result) {
		t.Error(fmt.Sprintf("Expected:\n%s\n\nGot:\n%s\n", expected, result))
	}
}
