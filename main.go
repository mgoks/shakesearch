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

const COMPLETE_WORKS_FILE = "completeworks.txt"
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

type Searcher struct {
	CompleteWorks						string
	CompleteWorksLowerCase	string
	SuffixArray							*suffixarray.Index
	// Index (int) to title (string) mappings.
	TitleMap								*treemap.Map
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

func (searcher *Searcher) Search(query string) []string {
	charOffset := 250
	idxs := searcher.SuffixArray.Lookup([]byte(strings.ToLower(query)), -1)		// []int

	// TODO Some results are missing the query. Fix it.
	/* TODO Rewrite this file and static/app.js so that this function returns
	* separate arrays for results, titles, indices, etc. instead of appending all
	* to results array. */
	// TODO Mark the query in the search results page.
	results := []string{}
	for _, idx := range idxs {
		var context string
		_, value := searcher.TitleMap.Floor(idx)
		title := value.(string)
		// Add 2 because otherwise context will start with ". "
		senBegin := lastIndexBefore(searcher.CompleteWorks, '.', idx) + 2;
		if senBegin < 0 {
			context = searcher.CompleteWorks[idx - charOffset / 2 : idx + charOffset / 2]
		} else {
			context = searcher.CompleteWorks[senBegin : senBegin + charOffset] + "..."
		}
		results = append(results, title, context)
	}
	return results
}

/* Returns the last index of a given character that comes before an index.
str: The string to search.
char: The character, the index of which this function returns.
index: Only the left of this index is searched.
*/
func lastIndexBefore(str string, char byte, index int) int {
	for i := index; i >= 0; i-- {
		if str[i] == char {
			return i
		}
	}
	return -1;
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
