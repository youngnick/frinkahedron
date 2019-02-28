package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nlopes/slack"
	"github.com/youngnick/frinkahedron/pkg/slackbot"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version    = "undefined"
	kingpinApp = kingpin.New("frinkahedron", "Slack bot to look up Frinkac things").DefaultEnvars().Version(version)

	verificationToken = kingpinApp.Arg("token", "The slack token").Default("").String()
)

func main() {
	kingpinApp.VersionFlag.Short('v')
	kingpinApp.HelpFlag.Short('h')
	kingpin.MustParse(kingpinApp.Parse(os.Args[1:]))

	http.HandleFunc("/slash", func(w http.ResponseWriter, r *http.Request) {
		s, err := slack.SlashCommandParse(r)
		fmt.Printf("%+v", s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// if !s.ValidateToken(verificationToken) {
		// 	fmt.Printf("%v", verificationToken)

		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		switch s.Command {
		case "/echo":
			params := &slack.Msg{Text: s.Text}
			b, err := json.Marshal(params)
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
	})

	http.HandleFunc("/slash/simpsons", slackbot.Simpsons)

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}
