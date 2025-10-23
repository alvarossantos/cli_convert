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
	delimiterFlag := flag.String("delimiter", ",", "csv delimiter")

	// Processa as flags e atribui os comandos as variaveis, obrigatório o uso do --
	flag.Parse()

	if *input == "" {
		fmt.Println("Missing required --input file")
		os.Exit(1)
	}

	if *output == "" {
		fmt.Println("Missing required --output file")
		os.Exit(1)
	}

	if *from == "" {
		fmt.Println("Missing required --from format")
		os.Exit(1)
	}

	if *to == "" {
		fmt.Println("Missing required --to format")
		os.Exit(1)
	}

	switch *from {
	case "json":
		if err := validateFileJSON(*input); err != nil {
			fmt.Printf("Error validating JSON: %v\n", err)
			os.Exit(1)
		}
	case "csv":
		runeArray := []rune(*delimiterFlag)
		if len(runeArray) != 1 {
			fmt.Println("Delimiter must be a single character")
			os.Exit(1)
		}
		if err := validateFileCSV(*input, runeArray[0]); err != nil {
			fmt.Printf("Error validating CSV: %v\n", err)
			os.Exit(1)
		}
	case "xml":
		if err := validateFileXML(*input); err != nil {
			fmt.Printf("Error validating XML: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unsupported format: %s\n", *from)
		os.Exit(1)
	}

	// Exibe os valores das flags
	fmt.Printf("Input: %s\n", *input)
	fmt.Printf("Output: %s\n", *output)
	fmt.Printf("From: %s\n", *from)
	fmt.Printf("To: %s\n", *to)
}
