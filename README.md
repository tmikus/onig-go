# Oniguruma Go

Go bindings for the Oniguruma regex library, a powerful and mature regular expression library with support for a wide range of character sets and language syntaxes. Oniguruma is written in C.

## Installation

**Prerequisites:**

In order to install onig-go, you need to have the Oniguruma library installed on your system. You can install it using Homebrew:

```bash
brew install oniguruma
```

**Installation:**

To install onig-go, use the following command:
```bash
go get github.com/tmikus/onig-go
```


## Example Usage

```go
package main

import (
    "fmt"
    "github.com/tmikus/onig-go"
)

func main() {
    regex, _ := onig.NewRegex("e(l+)")
    captures, _ := regex.Captures("hello")
    for _, text := range captures.All() {
        fmt.Println(text)
    }
}
```

## Documentation

The API documentation is available at [https://pkg.go.dev/github.com/tmikus/onig-go](https://pkg.go.dev/github.com/tmikus/onig-go).
