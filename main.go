package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"cli-convert/ai"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega .env se existir (não sobrescreve variáveis já definidas)
	godotenv.Load()

	if len(os.Args) < 2 {
		setGlobalUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--help", "-h":
		setGlobalUsage()
		os.Exit(0)

	case "convert":
		runConvert()

	case "detect":
		runDetect()

	case "schema":
		runSchema()

	case "ask":
		runAsk()

	default:
		fmt.Printf("Comando desconhecido: %s\n\n", os.Args[1])
		setGlobalUsage()
		os.Exit(1)
	}
}

// ──────────────────────────────────────────────
//  Comando: convert
// ──────────────────────────────────────────────

func runConvert() {
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)

	input := convertCmd.String("input", "", "arquivo de entrada")
	output := convertCmd.String("output", "", "arquivo de saída")
	from := convertCmd.String("from", "", "formato de origem (json, csv, xml, yaml)")
	to := convertCmd.String("to", "", "formato de destino (json, csv, xml, yaml)")
	delimiterFlag := convertCmd.String("delimiter", ",", "delimitador CSV")
	root := convertCmd.String("root", "root", "nome do elemento raiz para XML")
	convertCmd.Bool("help", false, "Mostra ajuda")

	setConvertUsage(convertCmd)

	// Verifica --help antes de parsear
	for _, arg := range os.Args[2:] {
		if arg == "--help" || arg == "-h" {
			convertCmd.Usage()
			os.Exit(0)
		}
	}

	convertCmd.Parse(os.Args[2:])

	// Validações
	if *input == "" {
		fmt.Println("Missing required --input file")
		os.Exit(1)
	}
	if *output == "" {
		fmt.Println("Missing required --output file")
		os.Exit(1)
	}

	// Auto-detecta formato se --from não foi especificado
	if *from == "" {
		detected, err := ai.DetectFormat(*input)
		if err != nil {
			fmt.Printf("Could not auto-detect format. Please specify --from.\nError: %v\n", err)
			os.Exit(1)
		}
		*from = detected
		fmt.Printf("Auto-detected format: %s\n", *from)
	}

	if *to == "" {
		fmt.Println("Missing required --to format")
		os.Exit(1)
	}

	// Valida formato de destino
	switch *to {
	case "json", "csv", "xml", "yaml":
		*output = ensureOutputExtension(*output, *to)
	default:
		fmt.Printf("Unsupported format: %s\n", *to)
		os.Exit(1)
	}

	// Abre arquivos
	fileIn, err := os.Open(*input)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer fileIn.Close()

	fileOut, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer fileOut.Close()

	// Valida delimitador
	runeArray := []rune(*delimiterFlag)
	if len(runeArray) != 1 {
		fmt.Println("Delimiter must be a single character")
		os.Exit(1)
	}
	delimiter := runeArray[0]

	// Dispatch de conversão
	if err := dispatchConversion(*from, *to, fileIn, fileOut, delimiter, *root); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Conversion from %s to %s completed successfully.\n", strings.ToUpper(*from), strings.ToUpper(*to))
}

// ──────────────────────────────────────────────
//  Comando: detect
// ──────────────────────────────────────────────

func runDetect() {
	detectCmd := flag.NewFlagSet("detect", flag.ExitOnError)
	input := detectCmd.String("input", "", "arquivo para detectar formato")
	detectCmd.Bool("help", false, "Mostra ajuda")

	detectCmd.Usage = func() {
		fmt.Println("cli-convert detect — Auto-detecta o formato de um arquivo.")
		fmt.Println()
		fmt.Println("USAGE:")
		fmt.Println("  cli-convert detect --input <file>")
		fmt.Println()
		fmt.Println("FLAGS:")
		fmt.Printf("  --input <string>   Arquivo para analisar\n")
		fmt.Printf("  -h, --help         Mostra esta ajuda\n")
	}

	for _, arg := range os.Args[2:] {
		if arg == "--help" || arg == "-h" {
			detectCmd.Usage()
			os.Exit(0)
		}
	}

	detectCmd.Parse(os.Args[2:])

	if *input == "" {
		fmt.Println("Missing required --input file")
		os.Exit(1)
	}

	format, err := ai.DetectFormat(*input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Detected format: %s\n", format)
}

// ──────────────────────────────────────────────
//  Comando: schema
// ──────────────────────────────────────────────

func runSchema() {
	schemaCmd := flag.NewFlagSet("schema", flag.ExitOnError)
	input := schemaCmd.String("input", "", "arquivo para inferir schema")
	schemaCmd.Bool("help", false, "Mostra ajuda")

	schemaCmd.Usage = func() {
		fmt.Println("cli-convert schema — Gera um JSON Schema a partir de um arquivo de dados.")
		fmt.Println()
		fmt.Println("USAGE:")
		fmt.Println("  cli-convert schema --input <file>")
		fmt.Println()
		fmt.Println("FLAGS:")
		fmt.Printf("  --input <string>   Arquivo para analisar\n")
		fmt.Printf("  -h, --help         Mostra esta ajuda\n")
		fmt.Println()
		fmt.Println("Para arquivos JSON, o schema é gerado localmente.")
		fmt.Println("Para outros formatos, usa IA (configure OPENROUTER_API_KEY no .env).")
	}

	for _, arg := range os.Args[2:] {
		if arg == "--help" || arg == "-h" {
			schemaCmd.Usage()
			os.Exit(0)
		}
	}

	schemaCmd.Parse(os.Args[2:])

	if *input == "" {
		fmt.Println("Missing required --input file")
		os.Exit(1)
	}

	schema, err := ai.InferSchema(*input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Imprime o schema formatado
	fmt.Printf("Schema for: %s\n", schema["source"])
	fmt.Printf("Format: %s\n\n", schema["format"])

	if props, ok := schema["properties"].(map[string]interface{}); ok && len(props) > 0 {
		fmt.Println("Properties:")
		for key, prop := range props {
			if p, ok := prop.(map[string]interface{}); ok {
				fmt.Printf("  %s: %v\n", key, p["type"])
			}
		}
	}

	if analysis, ok := schema["ai_analysis"].(string); ok {
		fmt.Printf("\nAI Analysis:\n%s\n", analysis)
	}
}

// ──────────────────────────────────────────────
//  Comando: ask
// ──────────────────────────────────────────────

func runAsk() {
	askCmd := flag.NewFlagSet("ask", flag.ExitOnError)
	input := askCmd.String("input", "", "arquivo para analisar")
	question := askCmd.String("question", "", "pergunta em linguagem natural")
	askCmd.Bool("help", false, "Mostra ajuda")

	askCmd.Usage = func() {
		fmt.Println("cli-convert ask — Pergunte sobre os dados de um arquivo em linguagem natural.")
		fmt.Println()
		fmt.Println("USAGE:")
		fmt.Println("  cli-convert ask --input <file> --question \"sua pergunta\"")
		fmt.Println()
		fmt.Println("FLAGS:")
		fmt.Printf("  --input <string>     Arquivo para analisar\n")
		fmt.Printf("  --question <string>  Pergunta em linguagem natural\n")
		fmt.Printf("  -h, --help           Mostra esta ajuda\n")
		fmt.Println()
		fmt.Println("Requer IA configurada (OPENROUTER_API_KEY no .env).")
	}

	for _, arg := range os.Args[2:] {
		if arg == "--help" || arg == "-h" {
			askCmd.Usage()
			os.Exit(0)
		}
	}

	askCmd.Parse(os.Args[2:])

	if *input == "" {
		fmt.Println("Missing required --input file")
		os.Exit(1)
	}
	if *question == "" {
		fmt.Println("Missing required --question")
		os.Exit(1)
	}

 resposta, err := ai.AskQuestion(*input, *question)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(resposta)
}
