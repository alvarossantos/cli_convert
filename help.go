package main

import (
	"flag"
	"fmt"
)

// Códigos ANSI para cores e estilos
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGray   = "\033[90m"
)

func setGlobalUsage() {
	fmt.Printf("%scli-convert%s — Conversor universal de arquivos JSON, CSV, XML e YAML.\n", ColorBold, ColorReset)
	fmt.Println()

	fmt.Printf("%sCOMANDOS DISPONÍVEIS%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %sconvert%s    Converte um arquivo entre formatos suportados\n", ColorYellow, ColorReset)
	fmt.Printf("  %sdetect%s     Auto-detecta o formato de um arquivo\n", ColorYellow, ColorReset)
	fmt.Printf("  %sschema%s     Gera um JSON Schema a partir de um arquivo\n", ColorYellow, ColorReset)
	fmt.Printf("  %sask%s        Pergunta sobre os dados em linguagem natural (IA)\n", ColorYellow, ColorReset)
	fmt.Println()

	fmt.Printf("%sEXEMPLOS%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %scli-convert convert --from json --to csv --input data.json --output data.csv%s\n", ColorGray, ColorReset)
	fmt.Printf("  %scli-convert detect --input arquivo.json%s\n", ColorGray, ColorReset)
	fmt.Printf("  %scli-convert schema --input dados.csv%s\n", ColorGray, ColorReset)
	fmt.Printf("  %scli-convert ask --input vendas.csv --question \"Qual o total de vendas?\"%s\n", ColorGray, ColorReset)
	fmt.Println()

	fmt.Printf("%sMAIS INFORMAÇÕES%s\n", ColorCyan, ColorReset)
	fmt.Println("  Execute 'cli-convert <comando> --help' para ajuda detalhada de cada comando.")
	fmt.Println()
}

func setConvertUsage(flagSet *flag.FlagSet) {
	flagSet.Usage = func() {
		fmt.Printf("%scli-convert convert%s — Converte um arquivo de um formato para outro.\n", ColorBold, ColorReset)
		fmt.Println()

		fmt.Printf("%sUSO%s\n", ColorCyan, ColorReset)
		fmt.Println("  cli-convert convert --from <formato> --to <formato> --input <arquivo> --output <arquivo> [flags]")
		fmt.Println()

		fmt.Printf("%sDESCRIÇÃO%s\n", ColorCyan, ColorReset)
		fmt.Println("  Converte arquivos entre JSON, CSV, XML e YAML, com opções de personalização")
		fmt.Println("  para delimitadores CSV e elementos raiz XML.")
		fmt.Println("  Se --from não for especificado, o formato é detectado automaticamente.")
		fmt.Println()

		fmt.Printf("%sFORMATOS SUPORTADOS%s\n", ColorCyan, ColorReset)
		fmt.Printf("  %s• json%s\n", ColorBlue, ColorReset)
		fmt.Printf("  %s• csv%s\n", ColorBlue, ColorReset)
		fmt.Printf("  %s• xml%s\n", ColorBlue, ColorReset)
		fmt.Printf("  %s• yaml%s\n", ColorBlue, ColorReset)
		fmt.Println()

		fmt.Printf("%sFLAGS OBRIGATÓRIOS%s\n", ColorCyan, ColorReset)
		fmt.Printf("  %s--from%s <string>       Formato de origem (json, csv, xml, yaml)\n", ColorYellow, ColorReset)
		fmt.Println("                       Opcional: detectado automaticamente se omitido")
		fmt.Println()
		fmt.Printf("  %s--to%s <string>         Formato de destino (json, csv, xml, yaml)\n", ColorYellow, ColorReset)
		fmt.Printf("  %s--input%s <string>      Caminho do arquivo de entrada\n", ColorYellow, ColorReset)
		fmt.Printf("  %s--output%s <string>     Caminho do arquivo de saída\n", ColorYellow, ColorReset)
		fmt.Println()

		fmt.Printf("%sFLAGS OPCIONAIS%s\n", ColorCyan, ColorReset)
		fmt.Printf("  %s--delimiter%s <char>    Delimitador CSV\n", ColorYellow, ColorReset)
		fmt.Println("                       Padrão: ','")
		fmt.Println()
		fmt.Printf("  %s--root%s <string>       Nome do elemento raiz para XML\n", ColorYellow, ColorReset)
		fmt.Println("                       Padrão: 'root'")
		fmt.Println()
		fmt.Printf("  %s-h%s, %s--help%s            Mostra esta mensagem de ajuda\n", ColorYellow, ColorReset, ColorYellow, ColorReset)
		fmt.Println()

		fmt.Printf("%sEXEMPLOS%s\n", ColorCyan, ColorReset)
		fmt.Printf("  %s# Converter JSON para CSV%s\n", ColorGray, ColorReset)
		fmt.Println("  cli-convert convert --from json --to csv --input data.json --output data.csv")
		fmt.Println()
		fmt.Printf("  %s# Converter CSV com delimitador ponto-e-vírgula para XML%s\n", ColorGray, ColorReset)
		fmt.Println("  cli-convert convert --from csv --to xml --input users.csv --output users.xml --delimiter ';' --root users")
		fmt.Println()
		fmt.Printf("  %s# Auto-detectar formato e converter para JSON%s\n", ColorGray, ColorReset)
		fmt.Println("  cli-convert convert --to json --input dados.csv --output dados.json")
		fmt.Println()
	}
}
