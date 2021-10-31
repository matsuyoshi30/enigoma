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
	fmt.Println(text)
	fmt.Println("===============   Cipher   ===============")
	c := e.Encrypt(text)
	fmt.Println(c)
	fmt.Println("===============   Decrypt  ===============")
	fmt.Println(e.Decrypt(c))
}
