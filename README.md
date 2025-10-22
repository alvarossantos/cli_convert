
# CLI Convert — Command-Line File Validator and Converter

**CLI Convert** is a tool written in **Go (Golang)** that allows you to validate and process files directly from the command line, with planned support for converting between multiple data formats (such as JSON → CSV, XML, YAML, etc.).

Currently, the project validates the structure of the input file and displays the provided parameters, serving as a foundation for future format conversion features.

---

## Features

* **JSON File Validation**

  Checks whether the file exists, is not empty, and contains valid JSON data.
* **Command-Line Interface (CLI)**

  Supports the parameters `--input`, `--output`, `--from`, and `--to`.
* **Extensible Architecture**

  Modular design makes it easy to add new format conversion functions in the future.

---

## 💻 Example Usage

```
Terminal:

./cli_convert --input input.json --output output.csv --from json --to csv
```
