package analyser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/martinrue/frekvenco/transcript"
)

type wordFrequencyMap map[string]int
type wordList []string

// WordFrequency holds a word and its frequency.
type WordFrequency struct {
	Word      string `json:"word"`
	Frequency int    `json:"frequency"`
}

// WordLength holds a word and its length.
type WordLength struct {
	Word   string `json:"word"`
	Length int    `json:"length"`
}

// Analysis holds a collection of stats about all words found in transcripts.
type Analysis struct {
	Transcripts int             `json:"transcripts"`
	Sentences   int             `json:"sentences"`
	Words       int             `json:"words"`
	Pronouns    []WordFrequency `json:"pronouns"`
	Tenses      []WordFrequency `json:"tenses"`
	Verbs       []WordFrequency `json:"verbs"`
	Compounds   []WordFrequency `json:"compounds"`
	Adverbs     []WordFrequency `json:"adverbs"`
	Nouns       []WordFrequency `json:"nouns"`
	Adjectives  []WordFrequency `json:"adjectives"`
	Longest     []WordLength    `json:"longest"`
	Top25       []WordFrequency `json:"top25"`
	Top100      []WordFrequency `json:"top100"`
	Top500      []WordFrequency `json:"top500"`
}

// Run builds language usage stats by analysing transcripts.
func Run(transcripts []*transcript.Transcript) (*Analysis, error) {
	wordFrequencies := make(wordFrequencyMap, 0)
	allWords := make(wordList, 0)

	analysis := &Analysis{}

	for _, transcript := range transcripts {
		for _, word := range transcript.Words {
			wordFrequencies[word]++
			allWords = append(allWords, word)
		}

		if len(transcript.Words) > 0 {
			analysis.Transcripts++
			analysis.Sentences += transcript.Lines
			analysis.Words += len(transcript.Words)
		}
	}

	top25, top100, top500 := calculateFrequencyMetrics(wordFrequencies)
	analysis.Pronouns = calculatePronounMetrics(wordFrequencies)
	analysis.Tenses = calculateTenseMetrics(wordFrequencies)
	analysis.Verbs = calculateVerbMetrics(wordFrequencies)
	analysis.Compounds = calculateCompoundVerbMetrics(allWords)
	analysis.Nouns = calculateNounMetrics(wordFrequencies)
	analysis.Adverbs = calculateAdverbMetrics(wordFrequencies)
	analysis.Adjectives = calculateAdjectiveMetrics(wordFrequencies)
	analysis.Longest = calculateLengthMetrics(wordFrequencies)
	analysis.Top25 = top25
	analysis.Top100 = top100
	analysis.Top500 = top500

	return analysis, nil
}

func calculateFrequencyMetrics(words wordFrequencyMap) (top25, top100, top500 []WordFrequency) {
	byFrequency := make([]WordFrequency, 0)

	for word, frequency := range words {
		byFrequency = append(byFrequency, WordFrequency{word, frequency})
	}

	sort.Slice(byFrequency, func(i, j int) bool {
		return byFrequency[i].Frequency > byFrequency[j].Frequency
	})

	return byFrequency[:25], byFrequency[25:125], byFrequency[125:625]
}

func calculatePronounMetrics(words wordFrequencyMap) []WordFrequency {
	forms := map[string]string{
		"mi":  "|mi|min|mia|mian|miaj|miajn|",
		"vi":  "|vi|vin|via|vian|viaj|viajn|",
		"li":  "|li|lin|lia|lian|liaj|liajn|",
		"ŝi":  "|ŝi|ŝin|ŝia|ŝian|ŝiaj|ŝiajn|",
		"ĝi":  "|ĝi|ĝin|ĝia|ĝian|ĝiaj|ĝiajn|",
		"ni":  "|ni|nin|nia|nian|niaj|niajn|",
		"ili": "|ili|ilin|ilia|ilian|iliaj|iliajn|",
		"oni": "|oni|onin|onia|onian|oniaj|oniajn|",
		"si":  "|si|sin|sia|sian|siaj|siajn|",
	}

	pronouns := wordFrequencyMap{
		"mi":  0,
		"vi":  0,
		"li":  0,
		"ŝi":  0,
		"ĝi":  0,
		"ni":  0,
		"ili": 0,
		"oni": 0,
		"si":  0,
	}

	for word, frequency := range words {
		for pronoun, variants := range forms {
			if strings.Contains(variants, fmt.Sprintf("|%s|", word)) {
				pronouns[pronoun] += frequency
			}
		}
	}

	byFrequency := make([]WordFrequency, 0)

	for pronoun, frequency := range pronouns {
		byFrequency = append(byFrequency, WordFrequency{pronoun, frequency})
	}

	sort.Slice(byFrequency, func(i, j int) bool {
		return byFrequency[i].Frequency > byFrequency[j].Frequency
	})

	return byFrequency
}

func calculateTenseMetrics(words wordFrequencyMap) []WordFrequency {
	tenses := wordFrequencyMap{
		"i":  0,
		"is": 0,
		"as": 0,
		"os": 0,
		"us": 0,
		"u":  0,
	}

	for word, frequency := range words {
		for tense := range tenses {
			if strings.HasSuffix(word, tense) {
				tenses[tense] += frequency
			}
		}
	}

	results := []WordFrequency{
		WordFrequency{"i", tenses["i"]},
		WordFrequency{"is", tenses["is"]},
		WordFrequency{"as", tenses["as"]},
		WordFrequency{"os", tenses["os"]},
		WordFrequency{"us", tenses["us"]},
		WordFrequency{"u", tenses["u"]},
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Frequency > results[j].Frequency
	})

	return results
}

func calculateVerbMetrics(words wordFrequencyMap) []WordFrequency {
	verbs := make([]WordFrequency, 0)

	for word, frequency := range words {
		endsWell := strings.HasSuffix(word, "i") ||
			strings.HasSuffix(word, "is") ||
			strings.HasSuffix(word, "as") ||
			strings.HasSuffix(word, "os") ||
			strings.HasSuffix(word, "us") ||
			strings.HasSuffix(word, "u")

		if len(word) >= 4 && endsWell {
			verbs = append(verbs, WordFrequency{word, frequency})
		}
	}

	sort.Slice(verbs, func(i, j int) bool {
		return verbs[i].Frequency > verbs[j].Frequency
	})

	return verbs[:25]
}

func calculateCompoundVerbMetrics(words wordList) []WordFrequency {
	compoundStart := func(word string) bool {
		return strings.HasSuffix(word, "is") ||
			strings.HasSuffix(word, "as") ||
			strings.HasSuffix(word, "os") ||
			strings.HasSuffix(word, "us")
	}

	compoundEnd := func(word string) bool {
		return strings.HasSuffix(word, "i")
	}

	compoundEndExceptions := map[string]bool{
		"mi":  true,
		"vi":  true,
		"ŝi":  true,
		"li":  true,
		"ĝi":  true,
		"ni":  true,
		"ili": true,
		"oni": true,
		"si":  true,
		"ci":  true,
		"pri": true,
		"pli": true,
		"ĉi":  true,
	}

	compounds := make(wordFrequencyMap, 0)

	for i := range words {
		if i == len(words)-1 {
			break
		}

		first := words[i]
		second := words[i+1]

		if compoundStart(first) && compoundEnd(second) {
			if _, ok := compoundEndExceptions[second]; !ok {
				compounds[fmt.Sprintf("%s %s", first, second)]++
			}
		}
	}

	byFrequency := make([]WordFrequency, 0)

	for compound, frequency := range compounds {
		byFrequency = append(byFrequency, WordFrequency{compound, frequency})
	}

	sort.Slice(byFrequency, func(i, j int) bool {
		return byFrequency[i].Frequency > byFrequency[j].Frequency
	})

	return byFrequency[:25]
}

func calculateNounMetrics(words wordFrequencyMap) []WordFrequency {
	nouns := make([]WordFrequency, 0)

	for word, frequency := range words {
		endsWell := strings.HasSuffix(word, "o") ||
			strings.HasSuffix(word, "on") ||
			strings.HasSuffix(word, "oj") ||
			strings.HasSuffix(word, "ojn")

		if len(word) >= 4 && endsWell {
			nouns = append(nouns, WordFrequency{word, frequency})
		}
	}

	sort.Slice(nouns, func(i, j int) bool {
		return nouns[i].Frequency > nouns[j].Frequency
	})

	return nouns[:25]
}

func calculateAdverbMetrics(words wordFrequencyMap) []WordFrequency {
	adverbs := make([]WordFrequency, 0)

	for word, frequency := range words {
		if len(word) >= 4 && strings.HasSuffix(word, "e") {
			adverbs = append(adverbs, WordFrequency{word, frequency})
		}
	}

	sort.Slice(adverbs, func(i, j int) bool {
		return adverbs[i].Frequency > adverbs[j].Frequency
	})

	return adverbs[:25]
}

func calculateAdjectiveMetrics(words wordFrequencyMap) []WordFrequency {
	adjectives := make([]WordFrequency, 0)

	pronouns := `|mi|min|mia|mian|miaj|miajn|
							 |vi|vin|via|vian|viaj|viajn|
							 |li|lin|lia|lian|liaj|liajn|
							 |ŝi|ŝin|ŝia|ŝian|ŝiaj|ŝiajn|
							 |ĝi|ĝin|ĝia|ĝian|ĝiaj|ĝiajn|
							 |ni|nin|nia|nian|niaj|niajn|
							 |ili|ilin|ilia|ilian|iliaj|iliajn|
							 |oni|onin|onia|onian|oniaj|oniajn|
							 |si|sin|sia|sian|siaj|siajn|`

	for word, frequency := range words {
		endsWell := strings.HasSuffix(word, "a") ||
			strings.HasSuffix(word, "an") ||
			strings.HasSuffix(word, "aj") ||
			strings.HasSuffix(word, "ajn")

		if len(word) >= 4 && endsWell && !strings.Contains(pronouns, fmt.Sprintf("|%s|", word)) {
			adjectives = append(adjectives, WordFrequency{word, frequency})
		}
	}

	sort.Slice(adjectives, func(i, j int) bool {
		return adjectives[i].Frequency > adjectives[j].Frequency
	})

	return adjectives[:25]
}

func calculateLengthMetrics(words wordFrequencyMap) []WordLength {
	byLength := make([]WordLength, 0)

	for word := range words {
		byLength = append(byLength, WordLength{word, len(word)})
	}

	sort.Slice(byLength, func(i, j int) bool {
		return byLength[i].Length > byLength[j].Length
	})

	return byLength[:25]
}
