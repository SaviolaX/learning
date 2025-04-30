package main

import (
	"fmt"
)
const (
    spanish = "Spanish"
    french = "French"
    englishHelloPrefix = "Hello, "
    spanishHelloPrefix = "Hola, "
    franchHelloPrefix = "Bonjour, "
)


func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }
    
    return greetingPrefix(language) + name 
}

func greetingPrefix(language string) (prefix string) {
    switch language {
        case spanish:
            prefix = spanishHelloPrefix
        case french:
            prefix = franchHelloPrefix
        default:
            prefix = englishHelloPrefix
    }
    return 
}

func main() {
    fmt.Println(Hello("world", "Spanish"))
}
