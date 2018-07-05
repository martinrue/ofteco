package transcript

import (
	"encoding/xml"
	"html"
	"regexp"
	"strings"
)

type document struct {
	XMLName xml.Name `xml:"transcript"`
	Text    []text   `xml:"text"`
}

type text struct {
	XMLName xml.Name `xml:"text"`
	Data    string   `xml:",chardata"`
}

func (t *text) getWords() []string {
	words := make([]string, 0)

	if strings.Contains(t.Data, "Amara") {
		return words
	}

	filter, err := regexp.Compile("[^a-z-ĉĝĥĵŝŭ ]+")
	if err != nil {
		return words
	}

	line := strings.Replace(t.Data, "\n", " ", -1)
	lower := strings.ToLower(html.UnescapeString(line))
	filtered := filter.ReplaceAllString(lower, "")

	for _, word := range strings.Split(filtered, " ") {
		if strings.TrimSpace(word) != "" {
			words = append(words, word)
		}
	}

	return words
}

func parseXML(bytes []byte) (*document, error) {
	doc := &document{}

	if len(bytes) > 0 {
		if err := xml.Unmarshal(bytes, doc); err != nil {
			return nil, err
		}
	}

	return doc, nil
}
