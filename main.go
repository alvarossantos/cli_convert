package main


import (
	"flag"
	"fmt"
)

func main() {

	// Definição das flags de linha de comando
	input := flag.String("input", "", "entrada")
	output := flag.String("output", "", "saida")
	from := flag.String("from", "", "origin format")
	to := flag.String("to", "", "target format")

	// Processa as flags e atribui os comandos as variaveis, obrigatório o uso do --
	flag.Parse()

	// Exibe os valores das flags
	fmt.Printf("Input: %s\n", *input)
	fmt.Printf("Output: %s\n", *output)
	fmt.Printf("From: %s\n", *from)
	fmt.Printf("To: %s\n", *to)
}