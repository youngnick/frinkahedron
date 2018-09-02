package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/youngnick/frinkahedron/pkg/frinkiac/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version    = "undefined"
	kingpinApp = kingpin.New("frinkacli", "Search a Frinkiac API for ").DefaultEnvars().Version(version)

	quote              = kingpinApp.Arg("quote", "Quote to search for").Required().String()
	text               = kingpinApp.Arg("text", "Text to overlay on the image").Default("").String()
	frinkiac           = kingpinApp.Flag("frinkiac", "Send the query to Frinkiac").Short('f').Bool()
	morbotron          = kingpinApp.Flag("morbotron", "Send the query to Morbotron").Short('m').Bool()
	masterofallscience = kingpinApp.Flag("masterofallscience", "Send the query to Masterofallscience").Short('c').Bool()
)

func main() {
	kingpinApp.VersionFlag.Short('v')
	kingpinApp.HelpFlag.Short('h')
	kingpin.MustParse(kingpinApp.Parse(os.Args[1:]))

	var apitarget *api.Frinkomatic

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

	r, err := http.Get(apitarget.ImageURL(frames[0], *text))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()

	iTermImgCat(base64.StdEncoding.EncodeToString(body))

}

func iTermImgCat(imagecontents string) {
	encodedFilename := base64.StdEncoding.EncodeToString([]byte("Frinkiac Image"))
	fmt.Printf("\n\033]1337;File=%v;inline=1:%v\n\n", encodedFilename, imagecontents)
}

func getBase64Image(filename string) string {
	imgFile, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer imgFile.Close()

	fInfo, _ := imgFile.Stat()
	size := fInfo.Size()

	buf := make([]byte, size)

	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	return base64.StdEncoding.EncodeToString(buf)

}
