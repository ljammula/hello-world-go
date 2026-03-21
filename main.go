package main

import (
    "flag"
    "fmt"
)

func main() {
    lang := flag.String("lang", "en", "Language for greeting: en or te")
    flag.Parse()

    switch *lang {
    case "en":
        fmt.Println("Hello, World!")
    case "te":
        fmt.Println("హలో వరల్డ్!") // Telugu: "Hello World!"
    default:
        fmt.Println("Unsupported language. Use 'en' or 'te'.")
    }
}
