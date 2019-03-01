package slackbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nlopes/slack"
	"github.com/youngnick/frinkahedron/pkg/frinkiac/api"
	"github.com/youngnick/frinkahedron/pkg/frinkiac/cmdparser"
)

func newImageMessage(title string, imageURL string, alttext string) slack.Message {
	var message slack.Message

	titleBlock := slack.NewTextBlockObject("plain_text", title, false, false)
	messageBlock := slack.NewImageBlock(imageURL, alttext, "image", titleBlock)
	message = slack.AddBlockMessage(message, messageBlock)
	message.ResponseType = "in_channel"

	return message
}

// FetchAGif will fetch a GIF with the provided parsed command, and send it back
// as a Slack message to responseUrl
func FetchAGif(apitarget *api.Frinkomatic, parsedCommand cmdparser.CommandDetails, responseURL string) {

	frames, err := apitarget.Search(parsedCommand.Quote)
	if err != nil {
		fmt.Println(err)
	}

	if len(frames) == 0 {
		fmt.Printf("No results found for %s\n", parsedCommand.Quote)
	}

	offset, err := time.ParseDuration(parsedCommand.GifOffset)
	if err != nil {
		//do an error
	}
	length, err := time.ParseDuration(parsedCommand.GifLength)
	if err != nil {
		//do an error
	}

	// These are used for the context frames function, we need to figure out
	// how to populate them correctly.
	var before time.Duration
	var after time.Duration

	// how before and after need to be set depends on offset's value

	if offset < 0 {
		before = offset * -1
		after = length - offset
	} else if offset > 0 {
		before = 0
		after = length + offset
	} else {
		before = 0
		after = length
	}

	var contextframes []api.Frame
	contextframes, err = apitarget.ContextFrames(frames[0], before, after)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Frames:\n%+v", contextframes)
	gifurl := apitarget.GifURL(contextframes[0], contextframes[len(contextframes)-1], parsedCommand.OverlayText)

	fmt.Printf("Gif URL: %v", gifurl)

	message := newImageMessage(parsedCommand.Original, gifurl, parsedCommand.Quote)

	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", responseURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
}
