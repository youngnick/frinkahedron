package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
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

func (f *Frinkomatic) getFrames(getURL string) ([]Frame, error) {
	var frames []Frame

	r, err := http.Get(getURL)
	if err != nil {
		return frames, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return frames, err
	}
	if string(body) == "[]" {
		return frames, nil
	}
	json.Unmarshal(body, &frames)
	return frames, nil
}

// Search searches for a quote on Frinkiac
func (f *Frinkomatic) Search(query string) ([]Frame, error) {

	checkURL := fmt.Sprintf("%v/api/search?q=%v", f.BaseURL, url.QueryEscape(query))

	return f.getFrames(checkURL)
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

// GifURL returns an GIF URL given a frame, and optionally, text to overlay
func (f *Frinkomatic) GifURL(firstframe Frame, lastframe Frame, text string) string {

	param := ""
	if len(text) > 0 {
		base64Text := base64.StdEncoding.EncodeToString([]byte(text))
		param = fmt.Sprintf("?b64lines=%v", base64Text)
	}

	return fmt.Sprintf("%v/gif/%v/%v/%v.gif%v", f.BaseURL, firstframe.Episode, firstframe.Timestamp, lastframe.Timestamp, param)
}

// ContextFrames retrieves a slice from frames from around the timestamp of the supplied Frame
func (f *Frinkomatic) ContextFrames(frame Frame, before time.Duration, after time.Duration) ([]Frame, error) {

	// url = u'{base}/api/frames/{episode}/{ts}/{before}/{after}'

	contextURL := fmt.Sprintf("%v/api/frames/%v/%v/%v/%v", f.BaseURL, frame.Episode, frame.Timestamp, int(before.Seconds())*1000, int(after.Seconds())*1000)

	return f.getFrames(contextURL)

}
