package renderer

import (
	"html/template"

	"github.com/martinrue/frekvenco/analyser"
	"github.com/martinrue/frekvenco/assets"
)

type model struct {
	Title      template.HTML
	CSS        template.CSS
	Header1    template.HTML
	Header2    template.HTML
	Logo       template.URL
	LogoLink   template.URL
	BookIcon   template.URL
	Top25Icon  template.URL
	Top100Icon template.URL
	Top500Icon template.URL
	Analysis   *analyser.Analysis
}

// Render renders an analysis into a self-contained HTML document.
func Render(analysis *analyser.Analysis, title, header1, header2, logo, logoLink string) (string, error) {
	logoData, err := fetchImage(logo)
	if err != nil {
		return "", err
	}

	return renderTemplate("/app.html", &model{
		Title:      template.HTML(title),
		CSS:        template.CSS(getAssetString("/app.css")),
		Header1:    template.HTML(header1),
		Header2:    template.HTML(header2),
		Logo:       template.URL(logoData),
		LogoLink:   template.URL(logoLink),
		BookIcon:   template.URL(getAssetImage("/book.svg", "image/svg+xml")),
		Top25Icon:  template.URL(getAssetImage("/top-25.svg", "image/svg+xml")),
		Top100Icon: template.URL(getAssetImage("/top-100.svg", "image/svg+xml")),
		Top500Icon: template.URL(getAssetImage("/top-500.svg", "image/svg+xml")),
		Analysis:   analysis,
	})
}

func getAssetString(name string) string {
	asset, _ := assets.FSString(false, name)
	return asset
}

func getAssetImage(name string, contentType string) string {
	bytes, err := assets.FSByte(false, name)
	if err != nil {
		return ""
	}

	return inlineImage(bytes, contentType)
}
