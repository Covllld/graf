package main

import (
	"fmt"

	gochoice "github.com/TwiN/go-choice"
)

func main() {
	choice, index, err := gochoice.Pick(
		"What do you want to do?\nPick:",
		[]string{
			"show system perf !",
			"Connect to the test environment",
			"Update",
		},
		gochoice.OptionBackgroundColor(gochoice.Black),
		gochoice.OptionTextColor(gochoice.White),
		gochoice.OptionSelectedTextColor(gochoice.Red),
		gochoice.OptionSelectedTextBold(),
	)

	if err != nil {
		fmt.Println("You didn't select anything!")
	} else {
		fmt.Printf("You have selected: '%s', which is the index %d\n", choice, index)

	}
}
