package slackbot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/youngnick/frinkahedron/pkg/frinkiac/api"

	"github.com/nlopes/slack"
	"github.com/youngnick/frinkahedron/pkg/frinkiac/cmdparser"
)

// Simpsons is the HandlerFunc that handles /simpsons commands for the bot
func Simpsons(w http.ResponseWriter, r *http.Request) {

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apitarget := api.New("frinkiac", "https://www.frinkiac.com", 24)

	// if !s.ValidateToken(verificationToken) {
	// 	fmt.Printf("%v", verificationToken)

	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	switch s.Command {
	case "/simpsons":
		parsed, err := cmdparser.Command(s.Text)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var message slack.Message

		frames, err := apitarget.Search(parsed.Quote)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		titleBlock := slack.NewTextBlockObject("plain_text", s.Text, false, false)

		messageBlock := slack.NewImageBlock(apitarget.ImageURL(frames[0], parsed.OverlayText),
			parsed.Quote,
			"image",
			titleBlock)
		message = slack.AddBlockMessage(message, messageBlock)
		message.ResponseType = "in_channel"

		b, err := json.Marshal(message)
		fmt.Println(string(b))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
