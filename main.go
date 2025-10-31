package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		setGlobalUsage()
		os.Exit(1)
	}

	if os.Args[1] == "convert" {
	
		convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)

		input := convertCmd.String("input", "", "entrada")
		output := convertCmd.String("output", "", "saida")
		from := convertCmd.String("from", "", "origin format")
		to := convertCmd.String("to", "", "target format")
		delimiterFlag := convertCmd.String("delimiter", ",", "csv delimiter")
		root := convertCmd.String("root", "root", "root element name for XML output (default: 'root')")
		convertCmd.Bool("help", false, "Mostra essa mensagem de ajuda")

		setConvertUsage(convertCmd)

		if len(os.Args) > 2 {
			for _, arg := range os.Args[2:] {
				if arg == "-help"|| arg == "--help" || arg == "-h" {
					convertCmd.Usage()
					os.Exit(0)
				}
			}
		}

		if len(os.Args) > 1 && os.Args[1] == "convert" {
			convertCmd.Parse(os.Args[2:])
		}
		if len(os.Args) < 2 || os.Args[1] != "convert" {
			fmt.Println("Expected 'convert' command")
			os.Exit(1)
		}

		convertCmd.Parse(os.Args[2:])

		if *input == "" {
			fmt.Println("Missing required --input file")
			flag.Usage()
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

			switch *to {
			case "csv":
				runeArray := []rune(*delimiterFlag)
				if len(runeArray) != 1 {
					fmt.Println("Delimiter must be a single character")
					os.Exit(1)
				}

				if !strings.HasSuffix(*output, ".csv") {
					*output += ".csv"
				}

				if err := convertJsonToCsv(*input, *output, runeArray[0]); err != nil {
					fmt.Printf("Error converting JSON to CSV: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("Conversion from JSON to CSV completed successfully.")
				os.Exit(0)

			case "xml":
				if !strings.HasSuffix(*output, ".xml") {
					*output += ".xml"
				}

				if err := convertJsonToXml(*input, *output, *root); err != nil {
					fmt.Printf("Error converting JSON to XML: %v", err)
					os.Exit(1)
				}

				fmt.Println("Conversion from JSON to XML completed sucessfully.")
				os.Exit(0)

			default:
				fmt.Printf("Unssuported conversion from JSON to %s", *to)
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

			switch *to {
			case "json":
				if !strings.HasSuffix(*output, ".json") {
					*output += ".json"
				}

				if err := convertCsvToJson(*input, *output, runeArray[0]); err != nil {
					fmt.Printf("Error converting CSV to JSON: %v", err)
					os.Exit(1)
				}
				fmt.Println("Conversion from CSV to JSON completed sucessfully.")
				os.Exit(0)
			case "xml":
				if !strings.HasSuffix(*output, ".xml") {
					*output += ".xml"
				}

				if err := convertCsvToXml(*input, *output, runeArray[0], *root); err != nil {
					fmt.Printf("Error converting CSV to XML: %v", err)
					os.Exit(1)
				}
				fmt.Println("Conversion from CSV to XML completed sucessfully.")
				os.Exit(0)
			default:
				fmt.Printf("Unssuported conversion from JSON to %s", *to)
				os.Exit(1)
			}

		case "xml":
			if err := validateFileXML(*input); err != nil {
				fmt.Printf("Error validating XML: %v\n", err)
				os.Exit(1)
			}

			switch *to {
			case "json":
				if !strings.HasSuffix(*output, ".json") {
					*output += ".json"
				}

				if err := convertXmlToJson(*input, *output); err != nil {
					fmt.Printf("Error converting XML to JSON: %v", err)
					os.Exit(1)
				}

				fmt.Println("Conversion from XML to JSON completed sucessfully.")
				os.Exit(0)

			case "csv":
				runeArray := []rune(*delimiterFlag)
				if len(runeArray) != 1 {
					fmt.Println("Delimiter must be a single character")
					os.Exit(1)
				}

				if !strings.HasSuffix(*output, ".csv") {
					*output += ".csv"
				}

				if err := convertXmlToCsv(*input, *output, runeArray[0]); err != nil {
					fmt.Printf("Error converting XML to CSV: %v", err)
					os.Exit(1)
				}

				fmt.Println("Conversion from XML to CSV completed sucessfully.")
				os.Exit(0)

			default:
				fmt.Printf("Unsuported conversion from XML to %s", *to)
				os.Exit(1)
			}

		default:
			fmt.Printf("Unsupported format: %s\n", *from)
			os.Exit(1)
		}
	}
}
