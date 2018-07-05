package renderer

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchImage(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return inlineImage(bytes, response.Header.Get("Content-Type")), nil
}

func inlineImage(bytes []byte, contentType string) string {
	data := base64.StdEncoding.EncodeToString(bytes)
	return fmt.Sprintf("data:%s;base64,%s\n", contentType, data)
}
