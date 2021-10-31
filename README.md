# Enigoma - Enigma + Go

A toy implementation for [Enigma machine](https://en.wikipedia.org/wiki/Enigma_machine) written in Go.

## Usage

Execute `go run _example/main.go_`.

```go
package main

import (
	"fmt"

	"github.com/matsuyoshi30/enigoma"
)

func main() {
	text := "I love Go because Gopher is so cute!!"

	pb := enigoma.NewPlugBoard()
	pb.AddExchange('a', 'e')
	pb.AddExchange('k', 't')
	pb.AddExchange('v', 'w')
	pb.AddExchange('j', 'c')
	pb.AddExchange('q', 'h')
	pb.AddExchange('g', 'n')

	e := enigoma.NewEnigoma(nil, nil, nil, 'd', 'k', 'x', pb)
	fmt.Println("=============== Plain Text ===============")
	fmt.Println(text)         // => I love Go because Gopher is so cute!!
	fmt.Println("===============   Cipher   ===============")
	c := e.Encrypt(text)
	fmt.Println(c)            // => H WYFU FJ VWQOJOS QTNLOS AJ IB YHXK!!
	fmt.Println("===============   Decrypt  ===============")
	fmt.Println(e.Decrypt(c)) // => i love go because gopher is so cute!!
}
```

## Licence

[MIT](./LICENSE)

## Author

[matsuyoshi30](https://twitter.com/matsuyoshi30)
