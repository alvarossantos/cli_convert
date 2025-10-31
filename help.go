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
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorGray   = "\033[90m"
)

func setGlobalUsage() {
	fmt.Printf("%scli-convert%s — A universal file format converter for JSON, XML, and CSV.\n", ColorBold, ColorReset)
	fmt.Println()

	fmt.Printf("%sUSAGE%s\n", ColorCyan, ColorReset)
	fmt.Println("  cli-convert <command> [flags]")
	fmt.Println()

	fmt.Printf("%sAVAILABLE COMMANDS%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %sconvert%s   Convert a file between supported formats\n", ColorGreen, ColorReset)
	fmt.Printf("  %shelp%s      Show general or command-specific help\n", ColorGreen, ColorReset)
	fmt.Println()

	fmt.Printf("%sEXAMPLES%s\n", ColorCyan, ColorReset)
	fmt.Printf("  %scli-convert convert --from json --to csv --input data.json --output data.csv%s\n", ColorGray, ColorReset)
	fmt.Println()

	fmt.Printf("%sMORE INFO%s\n", ColorCyan, ColorReset)
	fmt.Println("  Run 'cli-convert convert --help' for detailed usage of the 'convert' command.")
	fmt.Println()
}

func setConvertUsage(flagSet *flag.FlagSet) {
	flagSet.Usage = func() {
		fmt.Printf("%scli-convert convert%s — Convert a file from one format to another.\n", ColorBold, ColorReset)
		fmt.Println()

		fmt.Printf("%sUSAGE%s\n", ColorCyan, ColorReset)
		fmt.Println("  cli-convert convert --from <format> --to <format> --input <file> --output <file> [flags]")
		fmt.Println()

		fmt.Printf("%sDESCRIPTION%s\n", ColorCyan, ColorReset)
		fmt.Println("  Converts files between JSON, CSV, and XML formats, with optional customization")
		fmt.Println("  for CSV delimiters and XML root elements.")
		fmt.Println()

		fmt.Printf("%sREQUIRED FLAGS%s\n", ColorCyan, ColorReset)
		fmt.Printf("  %s--from%s <string>       Source format of the input file (options: json, csv, xml)\n", ColorYellow, ColorReset)
		fmt.Printf("  %s--to%s <string>         Target format of the output file (options: json, csv, xml)\n", ColorYellow, ColorReset)
		fmt.Printf("  %s--input%s <string>      Path to the input file to be converted\n", ColorYellow, ColorReset)
		fmt.Printf("  %s--output%s <string>     Path where the converted output file will be saved\n", ColorYellow, ColorReset)
		fmt.Println()

		fmt.Printf("%sOPTIONAL FLAGS%s\n", ColorCyan, ColorReset)
		fmt.Printf("  %s--delimiter%s <char>    CSV delimiter character (used for reading and writing)\n", ColorYellow, ColorReset)
		fmt.Println("                       Default: ','")
		fmt.Println()
		fmt.Printf("  %s--root%s <string>       Root element name for XML output files\n", ColorYellow, ColorReset)
		fmt.Println("                       Default: 'root'")
		fmt.Println()
		fmt.Printf("  %s-h%s, %s--help%s            Show this help message for the 'convert' command\n", ColorYellow, ColorReset, ColorYellow, ColorReset)
		fmt.Println()

		fmt.Printf("%sEXAMPLES%s\n", ColorCyan, ColorReset)
		fmt.Printf("  %s# Convert JSON to CSV%s\n", ColorGray, ColorReset)
		fmt.Println("  cli-convert convert --from json --to csv --input data.json --output data.csv")
		fmt.Println()
		fmt.Printf("  %s# Convert a semicolon-delimited CSV to XML with a custom root element%s\n", ColorGray, ColorReset)
		fmt.Println("  cli-convert convert --from csv --to xml --input users.csv --output users.xml --delimiter ';' --root users")
		fmt.Println()
		fmt.Printf("  %s# Convert XML to JSON%s\n", ColorGray, ColorReset)
		fmt.Println("  cli-convert convert --from xml --to json --input config.xml --output config.json")
		fmt.Println()
	}
}
