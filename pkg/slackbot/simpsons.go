package slackbot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/youngnick/frinkahedron/pkg/frinkiac/api"

	"github.com/nlopes/slack"
	"github.com/youngnick/frinkahedron/pkg/frinkiac/cmdparser"
)

// StrideStyle is the HandlerFunc that handles Stride-stle slash commands for the bot
// ie /<showname> <quote> [/ <giflength> [<gifoffset>]][| <overlaytext [...| <overlaytext]]
func StrideStyle(w http.ResponseWriter, r *http.Request) {

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var apitarget *api.Frinkomatic

	switch s.Command {
	case "/simpsons":
		apitarget = api.New("frinkiac", "https://www.frinkiac.com", 24)
	case "/futurama":
		apitarget = api.New("morbotron", "https://www.morbotron.com", 24)
	case "/rickandmorty":
		apitarget = api.New("masterofallscience", "https://www.masterofallscience.com", 24)
	}
	parsed, err := cmdparser.Command(s.Text)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	frames, err := apitarget.Search(parsed.Quote)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if parsed.GifLength != "s" {
		// We've been asked for a gif
		// First, respond to slack in this goroutine, then start a new
		// goroutine to return the gif (it will take a while)

		fmt.Printf("")
		gifMessage := &slack.Msg{Text: fmt.Sprintf("Getting your gif for '%v'", parsed.Quote)}

		g, err := json.Marshal(gifMessage)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		go FetchAGif(apitarget, parsed, s.ResponseURL)
		w.Header().Set("Content-Type", "application/json")
		w.Write(g)
	}

	message := newImageMessage(parsed.Original, apitarget.ImageURL(frames[0], parsed.OverlayText), parsed.Quote)

	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
