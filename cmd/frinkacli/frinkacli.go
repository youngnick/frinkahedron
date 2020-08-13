package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/youngnick/frinkahedron/pkg/frinkiac/api"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version    = "undefined"
	kingpinApp = kingpin.New("frinkacli", "Search a Frinkiac API for ").DefaultEnvars().Version(version)

	quote              = kingpinApp.Arg("quote", "Quote to search for").Required().String()
	offset             = kingpinApp.Arg("offset", "Offset from the quote location").Default("0").Duration()
	length             = kingpinApp.Arg("length", "length of a gif").Default("3s").Duration()
	text               = kingpinApp.Arg("text", "Text to overlay on the image").Default("").String()
	iterm              = kingpinApp.Flag("iterm", "Print the image using iTerm escape printing").Short('i').Bool()
	frinkiac           = kingpinApp.Flag("frinkiac", "Send the query to Frinkiac").Short('f').Bool()
	morbotron          = kingpinApp.Flag("morbotron", "Send the query to Morbotron").Short('m').Bool()
	masterofallscience = kingpinApp.Flag("masterofallscience", "Send the query to Masterofallscience").Short('c').Bool()
	gifmode            = kingpinApp.Flag("gifmode", "Get a gif instead of an image").Short('g').Bool()
)

func main() {
	kingpinApp.VersionFlag.Short('v')
	kingpinApp.HelpFlag.Short('h')
	kingpin.MustParse(kingpinApp.Parse(os.Args[1:]))

	var apitarget *api.Frinkomatic

	// Attempt to autodetect if you're using iTerm
	value, ok := os.LookupEnv("LC_TERMINAL")
	if ok && value == "iTerm2" {
		*iterm = true
	}
	if *morbotron {
		apitarget = api.New("morbotron", "https://www.morbotron.com", 24)
	} else if *masterofallscience {
		apitarget = api.New("masterofallscience", "https://www.masterofallscience.com", 24)
	} else {
		apitarget = api.New("frinkiac", "https://www.frinkiac.com", 24)
	}

	frames, err := apitarget.Search(*quote)
	if err != nil {
		log.Fatal(err)
	}
	if len(frames) == 0 {
		fmt.Printf("No results found for %s\n", *quote)
		os.Exit(0)
	}

	var body []byte
	var r *http.Response

	if *gifmode {
		var contextframes []api.Frame
		contextframes, err = apitarget.ContextFrames(frames[0], *offset, *length)
		if err != nil {
			log.Fatal(err)
		}
		gifurl := apitarget.GifURL(contextframes[0], contextframes[len(contextframes)-1], *text)
		longClient := http.Client{
			Timeout: time.Duration(120 * time.Second),
		}
		r, err = longClient.Get(gifurl)
		if err != nil {
			log.Fatal(err)
		}
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()
		fmt.Printf("%v\n", gifurl)
		if *iterm {
			iTermImgCat(base64.StdEncoding.EncodeToString(body))
			fmt.Print("\n")
		}
	} else {
		r, err = http.Get(apitarget.ImageURL(frames[0], *text))
		if err != nil {
			log.Fatal(err)
		}
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()
		fmt.Println(apitarget.ImageURL(frames[0], *text))
		if *iterm {
			iTermImgCat(base64.StdEncoding.EncodeToString(body))
			fmt.Print("\n")
		}

	}

}

func iTermImgCat(imagecontents string) {
	fmt.Print("\n")
	fmt.Printf("\033]1337;File=inline=1:%v\a\n", imagecontents)
	fmt.Print("\n")
}
