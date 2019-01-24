package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Frinkomatic holds the details for an instance of a Frinkiac-style API
type Frinkomatic struct {
	Name              string
	BaseURL           string
	SubtitleWrapWidth int
}

// New constructs a new Frinkomatic instance.
func New(name string, url string, subtitleWrapWidth int) *Frinkomatic {
	return &Frinkomatic{
		Name:              name,
		BaseURL:           url,
		SubtitleWrapWidth: subtitleWrapWidth,
	}
}

//A Frame Holds data about a Frinkiac (or Morbotron) search result
type Frame struct {
	ID        int    `json:"Id"`
	Episode   string `json:"Episode"`
	Timestamp int    `json:"Timestamp"`
}

// Search searches for a quote on Frinkiac
func (f *Frinkomatic) Search(query string) ([]Frame, error) {

	var frames []Frame

	checkURL := fmt.Sprintf("%v/api/search?q=%v", f.BaseURL, url.QueryEscape(query))

	fmt.Println(checkURL)
	r, err := http.Get(checkURL)
	if err != nil {
		return frames, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return frames, err
	}
	if string(body) == "[]" {
		return frames, errors.New("No results found for this search")
	}
	json.Unmarshal(body, &frames)
	return frames, nil
}

// ImageURL returns an Image URL given a frame, and optionally, text to overlay
func (f *Frinkomatic) ImageURL(frame Frame, text string) string {

	param := ""
	if len(text) > 0 {
		base64Text := base64.StdEncoding.EncodeToString([]byte(text))
		param = fmt.Sprintf("?b64lines=%v", base64Text)
	}

	return fmt.Sprintf("%v/meme/%v/%v.jpg%v", f.BaseURL, frame.Episode, frame.Timestamp, param)
}
