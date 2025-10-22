package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	// Definição das flags de linha de comando
	input := flag.String("input", "", "entrada")
	output := flag.String("output", "", "saida")
	from := flag.String("from", "", "origin format")
	to := flag.String("to", "", "target format")

	// Processa as flags e atribui os comandos as variaveis, obrigatório o uso do --
	flag.Parse()

	if err := validateFileJSON(*input); err != nil {
		fmt.Printf("Error validating JSON: %v\n", err)
		os.Exit(1)
	}

	// Exibe os valores das flags
	fmt.Printf("Input: %s\n", *input)
	fmt.Printf("Output: %s\n", *output)
	fmt.Printf("From: %s\n", *from)
	fmt.Printf("To: %s\n", *to)
}
