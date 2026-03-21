package main

import (
    "flag"
    "fmt"
    "os"
)

var greetings = map[string]string{
    "en": "Hello, World!",
    "te": "హలో వరల్డ్!",
}

func main() {
    lang := flag.String("lang", "en", "Language for greeting: en or te")
    flag.Parse()

    if greet, ok := greetings[*lang]; ok {
        fmt.Println(greet)
    } else {
        fmt.Fprintln(os.Stderr, "Unsupported language. Use 'en' or 'te'.")
        os.Exit(1)
    }
}
