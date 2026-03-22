package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

var greetings = map[string]string{
	"en": "Hello, World!",
	"te": "హలో వరల్డ్!",
	"hi": "नमस्ते दुनिया!",
}

func supportedLangs() string {
	keys := make([]string, 0, len(greetings))
	for k := range greetings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func main() {
	usage := "Language for greeting: " + supportedLangs()
	lang := flag.String("lang", "en", usage)
	mode := flag.String("mode", "cli", "Execution mode: cli or http")
	addr := flag.String("addr", ":8080", "HTTP server address")
	indexFile := flag.String("index", "index.html", "Path to index file for HTTP mode")
	flag.Parse()

	switch *mode {
	case "http":
		if err := runHTTPServer(*addr, *indexFile); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
			os.Exit(1)
		}
	case "cli":
		if greet, ok := greetings[*lang]; ok {
			fmt.Println(greet)
		} else {
			fmt.Fprintf(os.Stderr, "Unsupported language. Use one of: %s.\n", supportedLangs())
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported mode %q. Use one of: cli, http.\n", *mode)
		os.Exit(1)
	}
}

func runHTTPServer(addr string, indexFile string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if _, err := os.Stat(indexFile); err != nil {
			if os.IsNotExist(err) {
				http.Error(w, fmt.Sprintf("index file not found: %s", indexFile), http.StatusNotFound)
				return
			}
			http.Error(w, "failed to read index file", http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, indexFile)
	})

	return http.ListenAndServe(addr, mux)
}
