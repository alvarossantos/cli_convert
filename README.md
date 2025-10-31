# CLI Convert – Universal File Converter

**CLI Convert** is a powerful and flexible command-line interface (CLI) tool built with **Go (Golang)**, designed to streamline the conversion of data files between various formats. It currently supports seamless transformations between **JSON, CSV, and XML**, with a focus on performance and ease of use.

This tool is ideal for developers, data engineers, and anyone who frequently works with different data formats and needs a quick, reliable way to convert them.

---

## ✨ Features

* **Bidirectional Conversion:** Convert between JSON, CSV, and XML formats.
  * `JSON <-> CSV`
  * `JSON <-> XML`
  * `CSV <-> XML`
* **Robust File Validation:** Ensures input files exist, are not empty, and adhere to their specified format before conversion.
* **Intuitive Command-Line Interface:** Easy-to-use flags for specifying input/output, source/target formats, and conversion-specific options.
* **Go-powered Performance:** Leverages Go's concurrency and efficiency for fast data processing.
* **Extensible Architecture:** Modular design allows for easy addition of new data formats and conversion logic in the future.
* **Intelligent Type Handling:** Automatically detects and converts numerical and boolean values from source formats into appropriate JSON types, preventing numbers from being treated as strings.
* **XML to JSON Order Preservation:** When converting XML to JSON, the tool now preserves the order of elements, ensuring a more faithful representation of the original data structure.

---

## 🚀 Installation

To get started with CLI Convert, ensure you have Go installed on your system. Then, you can build the executable:

```bash
# Clone the repository (if not already done)
git clone https://github.com/your-username/cli_convert.git # Replace with actual repo URL
cd cli_convert

# Build the executable
go build

# (Optional) Move the executable to your PATH for global access
mv cli-convert /usr/local/bin/
```

---

## 📖 Usage

The `cli-convert` tool uses a single mandatory command: `convert`. All operations are performed using this command followed by specific flags.

### Basic Syntax

```bash
cli-convert convert --from <source_format> --to <target_format> --input <input_file> --output <output_file> [options]
```

### Available Flags

* `--input <file_path>` (Required)

  * Specifies the path to the source file you want to convert.
  * Example: `--input data.json`
* `--output <file_path>` (Required)

  * Specifies the path where the converted file will be saved.
  * The tool will automatically append the correct file extension (`.json`, `.csv`, `.xml`) if not provided.
  * Example: `--output converted_data.csv`
* `--from <format>` (Required)

  * Defines the format of the input file.
  * Accepted values: `json`, `csv`, `xml`
  * Example: `--from json`
* `--to <format>` (Required)

  * Defines the desired format for the output file.
  * Accepted values: `json`, `csv`, `xml`
  * Example: `--to csv`
* `--delimiter <char>` (Optional)

  * Used when converting to or from CSV files.
  * Specifies the character used to separate values in the CSV.
  * Default: `,` (comma)
  * Example: `--delimiter ';'` (for semicolon-separated CSV)
* `--root <string>` (Optional)

  * Used specifically when converting from JSON to XML.
  * Defines the name of the root element in the generated XML output.
  * Default: `root`
  * Example: `--root MyDataCollection`

### Conversion Examples

Here are some common conversion scenarios:

#### ➡️ JSON to CSV

```bash
cli-convert convert --from json --to csv --input input.json --output output.csv# Using a custom delimiter
cli-convert convert --from json --to csv --input input.json --output output.csv --delimiter ';'

#### ➡️ JSON to XML

```bash
cli-convert convert --from json --to xml --input input.json --output output.xml
# With a custom root element name
cli-convert convert --from json --to xml --input input.json --output output.xml --root MyJsonData
```

#### ➡️ CSV to JSON

```bash
cli-convert convert --from csv --to json --input input.csv --output output.json
# With a custom delimiter for the input CSV
cli-convert convert --from csv --to json --input input.csv --output output.json --delimiter ';'
```

#### ➡️ CSV to XML

```bash
cli-convert convert --from csv --to xml --input input.csv --output output.xml
# With a custom delimiter for the input CSV and a custom root element
cli-convert convert --from csv --to xml --input input.csv --output output.xml --delimiter ';' --root CsvRecords
```

#### ➡️ XML to JSON

```bash
cli-convert convert --from xml --to json --input input.xml --output output.json
```

#### ➡️ XML to CSV

```bash
cli-convert convert --from xml --to csv --input input.xml --output output.csv
cli-convert convert --from xml --to csv --input input.xml --output output.csv --delimiter '|'
```

---

## ⚠️ Error Handling

The tool provides informative error messages for common issues such as:

* Missing required flags (`--input`, `--output`, `--from`, `--to`).
* Invalid or non-existent input files.
* Unsupported conversion formats.
* Malformed input data (e.g., invalid JSON, CSV, or XML structure).
* Delimiter not being a single character.

Always check the output for error messages if a conversion fails.

---

## 🤝 Contributing

Contributions are welcome! If you have suggestions for new features, bug fixes, or improvements, please feel free to open an issue or submit a pull request.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. (Note: Create a LICENSE file if you haven't already.)
