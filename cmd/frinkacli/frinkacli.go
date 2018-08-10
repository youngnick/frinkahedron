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
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	quote   = kingpin.Arg("name", "Name of user.").Required().String()
)

func main() {
	kingpin.Parse()

	frinkiac := api.New("frinkiac", "https://www.frinkiac.com", 24)
	frames, err := frinkiac.Search(*quote)
	if err != nil {
		log.Fatal(err)
	}

	r, err := http.Get(frinkiac.ImageURL(frames[0]))
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
