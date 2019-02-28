// Package cmdparser provides bits to parse out queries from a old-school slash command
package cmdparser

import (
	"strings"
)

// CommandDetails records
type CommandDetails struct {
	Quote       string `json:"quote"`
	GifLength   string `json:"giflength"`
	GifOffset   string `json:"gifoffset"`
	OverlayText string `json:"overlaytext"`
}

// Command parses a standard command test and returns a CommandDetails struct with all the deets
// The format is <quote text> / <gif length> <gif offset> | <overlaytext line 1> | <overlay line 2> ...
func Command(commandtext string) (CommandDetails, error) {

	var parsedCommand CommandDetails

	// Let's tokenise this using the delimiters, that'll make this parsing easier.
	// super fragile parsing up ahead!
	overlaySlice := strings.Split(commandtext, "|")

	if len(overlaySlice) > 1 {
		parsedCommand.OverlayText = strings.Join(overlaySlice[1:], "\n")
		commandtext = overlaySlice[0]
	}

	gifSlice := strings.Split(commandtext, "/")

	if len(gifSlice) > 1 {
		gifDetails := strings.Split(strings.TrimSpace(gifSlice[1]), " ")
		if len(gifDetails) > 1 {
			parsedCommand.GifOffset = gifDetails[1]

		}
		parsedCommand.GifLength = gifDetails[0]
		commandtext = gifSlice[0]
	}

	parsedCommand.Quote = commandtext

	// fmt.Printf("%+v", result)

	return parsedCommand, nil
}
