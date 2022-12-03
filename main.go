package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Searcher struct {
	CompleteWorks          string
	CompleteWorksLowerCase string
	SuffixArray            *suffixarray.Index
	TitleMap               *treemap.Map
}

type Result struct {
	context string
	title   string
	begin   int
	end     int
}

const COMPLETE_WORKS_FILE = "completeworks.txt"
const PUNCT_MARKS = ".!?]"
const TITLES_FILE = "titles.txt"

func main() {
	searcher := Searcher{}
	err := searcher.Load(COMPLETE_WORKS_FILE)
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		// TODO Marshal results to JSON and send that to front-end instead of
		// sending results directly.
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (searcher *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	searcher.CompleteWorks = string(dat)
	searcher.CompleteWorksLowerCase = strings.ToLower(searcher.CompleteWorks)
	searcher.SuffixArray = suffixarray.New([]byte(searcher.CompleteWorksLowerCase))
	searcher.createTitleMap()
	return nil
}

// TODO Change return type to []Result
func (searcher *Searcher) Search(query string) []string {
	var idxs []int = searcher.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)
	offset := 250
	results := []string{}
	for _, queryBegin := range idxs { // index of query within searcher.CompleteWorks
		_, value := searcher.TitleMap.Floor(queryBegin)
		title := value.(string)
		// Begin and end indices of context within searcher.CompleteWorks.
		// Add 2 because otherwise context will start with ". "
		contextBegin := lastIndexBefore(searcher.CompleteWorks, PUNCT_MARKS, queryBegin) + 2
		contextEnd := contextBegin + offset
		// Make sure query is in context.
		if queryBegin+len(query) > contextEnd {
			contextEnd = queryBegin + len(query)
		}
		context := searcher.CompleteWorks[contextBegin:contextEnd] + "..."
		// Fall back to original method if no end of previous sentence found.
		if contextBegin < 0 {
			context = searcher.CompleteWorks[queryBegin-offset/2 : queryBegin+offset/2]
		}
		// Begin index of query within context.
		b := queryBegin - contextBegin
		e := b + len(query)
		// TODO Move marking to front-end.
		context = context[:b] + "<mark>" + context[b:e] + "</mark>" + context[e:]
		results = append(results, title, context)
	}
	return results
}

/*
	Returns the last index of given characters that comes before an index.

str: The string to search.
chars: The characters, the index of which this function returns.
index: Only the left of this index is searched.
*/
func lastIndexBefore(str string, chars string, index int) int {
	for i := index; i >= 0; i-- {
		if strings.IndexByte(chars, str[i]) >= 0 {
			return i
		}
	}
	return -1
}

/* Populates searcher's TitleMap such that the keys refer to the beginning of a
 * work in CompleteWorks string, the name of which is the value, e.g. assume the
 * Sonnets start at position 0 and Hamlet starts at 92380, then there will be 2
 * key-value pairs in the map such as 0 and 'The Sonnets', and 92380 and
 * 'Hamlet.'
 */
func (searcher *Searcher) createTitleMap() error {
	data, err := ioutil.ReadFile(TITLES_FILE)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	titles := strings.Split(string(data), "\n")
	searcher.TitleMap = treemap.NewWithIntComparator()
	for _, title := range titles {
		lastIndex := strings.LastIndex(searcher.CompleteWorks, title)
		searcher.TitleMap.Put(lastIndex, title)
	}
	return nil
}
