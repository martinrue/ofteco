package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/martinrue/frekvenco/analyser"
	"github.com/martinrue/frekvenco/renderer"
	"github.com/martinrue/frekvenco/transcript"
)

const usage = `Frekvenco (v0.0.1)

Usage:
  frekvenco [config]

Application Config:
  --videos=     path to input file containing video IDs
  --title=      page title in output document
  --header-1=   primary header in output document
  --header-2=   secondary header in output document
  --logo=       URL of logo image in output document
  --logo-link=  URL of logo link in output document
`

var (
	videos   = flag.String("videos", "", "")
	title    = flag.String("title", "", "")
	header1  = flag.String("header-1", "", "")
	header2  = flag.String("header-2", "", "")
	logo     = flag.String("logo", "", "")
	logoLink = flag.String("logo-link", "", "")
)

func validateFlags() {
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			flag.Usage()
			fmt.Fprintln(os.Stderr, "\nerror: all config flags must be specified")
			os.Exit(1)
		}
	})
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()

	validateFlags()

	ids, err := readInputFile(*videos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not read input file → %s\n", err)
		os.Exit(3)
	}

	analysis, err := analyser.Run(fetchTranscripts(ids))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed analysing transcripts → %s\n", err)
		os.Exit(4)
	}

	fmt.Fprintf(os.Stderr, "\n✓ analysed %d words from %d transcripts\n", analysis.Words, analysis.Transcripts)

	content, err := renderer.Render(analysis, *title, *header1, *header2, *logo, *logoLink)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed rendering → %s\n", err)
	}

	fmt.Print(content)
}

func readInputFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	ids := make([]string, 0)

	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}

	return ids, scanner.Err()
}

func fetchTranscripts(ids []string) []*transcript.Transcript {
	transcripts := make([]*transcript.Transcript, 0)

	for _, id := range ids {
		t, err := transcript.Fetch(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ignoring [%s] → %s\n", id, err)
			continue
		}

		fmt.Fprintf(os.Stderr, "fetched %d words from [%s]\n", len(t.Words), t.Video)
		transcripts = append(transcripts, t)
	}

	return transcripts
}
