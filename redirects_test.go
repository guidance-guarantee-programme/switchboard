package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

var responseXML = []string{
	`
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Dial>
        <Number>789</Number>
    </Dial>
</Response>
	`,
	`
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Dial>
        <Number sendDigits="100">789</Number>
    </Dial>
</Response>
	`,
}

var responseJSON = []byte(`
{
  "phone": "789"
}
`)

func TestLoadRedirectsFromYAML(t *testing.T) {
	var (
		redirects []Redirect
		redirect  Redirect
	)

	redirects = LoadRedirectsFromYAML("fixtures/redirects.yaml")

	if len(redirects) != 2 {
		t.Error("Should return a slice of 2 Redirects")
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

	if redirect.CabExtension != "" {
		t.Error("CabExtension should not be set")
	}

	redirectWithExtension := redirects[1]

	if redirectWithExtension.CabExtension != "100" {
		t.Error("CabExtension should be set")
	}
}

func TestFindRedirectForTwilio(t *testing.T) {
	redirects := LoadRedirectsFromYAML("fixtures/redirects.yaml")

	if _, err := FindRedirectForTwilio(redirects, "456"); err != nil {
		t.Error("Should not return an error for valid Twilio number")
	}

	if redirect, _ := FindRedirectForTwilio(redirects, "456"); redirect != redirects[0] {
		t.Error("Should return Redirect for valid Twilio number")
	}

	if _, err := FindRedirectForTwilio(redirects, "999"); err == nil {
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
	redirects := []Redirect{
		{Cab: "789"},
		{Cab: "789", CabExtension: "100"},
	}

	for i := 0; i <= len(responseXML)-1; i++ {
		expected := []byte(strings.TrimSpace(responseXML[i]))
		result := GenerateResponseXMLFor(redirects[i])

		if !bytes.Equal(expected, result) {
			t.Error(fmt.Sprintf("Expected:\n%s\n\nGot:\n%s\n", expected, result))
		}
	}
}

func TestGenerateResponseJSONFor(t *testing.T) {
	expected := new(bytes.Buffer)
	_ = json.Compact(expected, responseJSON)

	result := GenerateResponseJSONFor("789")

	if !bytes.Equal(expected.Bytes(), result) {
		t.Error(fmt.Sprintf("Expected:\n%s\n\nGot:\n%s\n", expected, result))
	}
}
