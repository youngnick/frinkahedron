package main

import (
	"fmt"
	"net/http"
	"os"

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

	http.HandleFunc("/stridestyle", slackbot.StrideStyle)

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8080", nil)
}
