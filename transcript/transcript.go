package transcript

import (
	"io/ioutil"
	"net/http"
)

// Transcript holds data about a video's Esperanto transcript.
type Transcript struct {
	Video string
	Lines int
	Words []string
}

// Fetch retrieves an Esperanto video transcription from a YouTube video ID.
func Fetch(id string) (*Transcript, error) {
	response, err := http.Get("https://video.google.com/timedtext?lang=eo&v=" + id)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	document, err := parseXML(bytes)
	if err != nil {
		return nil, err
	}

	transcript := &Transcript{Video: id, Lines: len(document.Text)}

	for _, text := range document.Text {
		for _, word := range text.getWords() {
			transcript.Words = append(transcript.Words, word)
		}
	}

	return transcript, nil
}
