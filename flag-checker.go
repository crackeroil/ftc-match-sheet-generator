package main

import (
	"flag"
	"fmt"
)

func CheckFlags() (bool, string, string, string) {
	var goodFlags bool = true

	var helpFlag bool
	var seasonFlag string
	var eventCodeFlag string
	var apiKeyFlag string

	flag.BoolVar(&helpFlag, "h", false, "Show this help message")
	flag.StringVar(&seasonFlag, "s", "", "The season of the event")
	flag.StringVar(&eventCodeFlag, "e", "", "The event code of the event")
	flag.StringVar(&apiKeyFlag, "k", "", "The API key to use")

	flag.Parse()

	if helpFlag {
		fmt.Println("flags:")

		flag.PrintDefaults()

		goodFlags = false
	}

	if !helpFlag && seasonFlag == "" {
		fmt.Println("Season flag is required")

		goodFlags = false
	}

	if !helpFlag && eventCodeFlag == "" {
		fmt.Println("Event code flag is required")

		goodFlags = false
	}

	if !helpFlag && apiKeyFlag == "" {
		fmt.Println("API key is required\n(Should be in the form of \"username:authorizationKey\")")

		goodFlags = false
	}

	return goodFlags, seasonFlag, eventCodeFlag, apiKeyFlag
}
