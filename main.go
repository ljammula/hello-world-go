package main

import (
    "flag"
    "fmt"
    "os"
)

var greetings = map[string]string{
    "en": "Hello, World!",
    "te": "హలో వరల్డ్!",
    "hi": "नमस्ते दुनिया!",
}

func main() {
    lang := flag.String("lang", "en", "Language for greeting: en, te, or hi")
    flag.Parse()

    if greet, ok := greetings[*lang]; ok {
        fmt.Println(greet)
    } else {
        fmt.Fprintln(os.Stderr, "Unsupported language. Use 'en', 'te', or 'hi'.")
        os.Exit(1)
    }
}
