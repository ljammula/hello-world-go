package main

import (
    "flag"
    "fmt"
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
    flag.Parse()

    if greet, ok := greetings[*lang]; ok {
        fmt.Println(greet)
    } else {
        fmt.Fprintf(os.Stderr, "Unsupported language. Use one of: %s.\n", supportedLangs())
        os.Exit(1)
    }
}
