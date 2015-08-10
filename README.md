# Switchboard

[![Build Status](https://travis-ci.org/guidance-guarantee-programme/switchboard.svg)](https://travis-ci.org/guidance-guarantee-programme/switchboard)

We're proxying the phone numbers used to contact the Citizens Advice service providing [Pension Wise] face to face guidance appointments through [Twilio] in order to provide extra anayltics.

This app provides a lookup for the [main Rails app](https://github.com/guidance-guarantee-programme/pension_guidance) and handles the webhooks from [Twilio] when a call is initiated.


## Prerequisites

* [Go]
* [Git]
* [A configured GOPATH](https://github.com/golang/go/wiki/GOPATH)


## Installation

Clone the repository:

```sh
$ cd $GOPATH/src
$ git clone https://github.com/guidance-guarantee-programme/switchboard.git
```

Install godep:

```sh
go get github.com/tools/godep
```

## Usage

To start the application:

```sh
$ godep go build .
$ ./switchboard
```

## Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


## Contributing

Please see the [contributing guidelines](/CONTRIBUTING.md).

[git]: http://git-scm.com
[go]: https://golang.org
[pension wise]: https://www.pensionwise.gov.uk
[twilio]: https://www.twilio.com/
