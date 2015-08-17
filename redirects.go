package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Redirect struct {
	Id           string `yaml:"id"`
	Twilio       string `yaml:"twilio"`
	Cab          string `yaml:"cab"`
	CabExtension string `yaml:"cab_ext"`
}

type Number struct {
	Number    string `xml:",chardata"`
	Extension string `xml:"sendDigits,attr,omitempty"`
}

type Dial struct {
	Number Number
}

type Response struct {
	Dial Dial
}

type Phone struct {
	Phone string `json:"phone"`
}

func LoadRedirectsFromYAML(path string) (redirects []Redirect) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Can't read %s", path))
	}

	err = yaml.Unmarshal(data, &redirects)
	if err != nil {
		panic(fmt.Sprintf("Can't parse %s", path))
	}

	return
}

func FindRedirectForTwilio(redirects []Redirect, twilio string) (Redirect, error) {
	for _, redirect := range redirects {
		if redirect.Twilio == twilio {
			return redirect, nil
		}
	}

	return Redirect{}, errors.New("Redirect not found")
}

func FindTwilioForID(redirects []Redirect, id string) (string, error) {
	for _, redirect := range redirects {
		if redirect.Id == id {
			return redirect.Twilio, nil
		}
	}

	return "", errors.New("Redirect not found")
}

func GenerateResponseXMLFor(redirect Redirect) []byte {
	response := &Response{
		Dial: Dial{
			Number: Number{
				Number:    redirect.Cab,
				Extension: redirect.CabExtension,
			},
		},
	}

	responseXML, err := xml.MarshalIndent(response, "", "    ")
	if err != nil {
		panic("Couldn't generate XML")
	}

	result := bytes.NewBufferString(xml.Header)
	result.Write(responseXML)

	return result.Bytes()
}

func GenerateResponseJSONFor(phone string) []byte {
	response, err := json.Marshal(&Phone{phone})

	if err != nil {
		panic("Can't Generate JSON")
	}

	return response
}
