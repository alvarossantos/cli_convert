# CLI Convert — Command-Line File Validator and Converter

**CLI Convert** is a tool written in **Go (Golang)** that allows you to validate and process files directly from the command line, with planned support for converting between multiple data formats (such as JSON → CSV, XML, YAML, etc.).

Currently, the project validates the structure of the input file and displays the provided parameters, serving as a foundation for future format conversion features.

---

## Features

* **File Validation**

  Checks whether the file exists, is not empty, and contains valid data.
* **Command-Line Interface (CLI)**

  Supports the parameters `--input`, `--output`, `--from`, and `--to`.
* **Extensible Architecture**

  Modular design makes it easy to add new format conversion functions in the future.

---


### 📄 Project Description

**CLI Convert** is a lightweight and extensible command-line utility built with **Go (Golang)** for validating and processing data files.

Its main purpose is to provide a simple interface to check the integrity of JSON files and prepare the foundation for future data format conversions such as  **JSON → CSV** ,  **JSON → XML** , or  **YAML → JSON** .

The tool includes a flexible CLI interface that accepts custom parameters for input and output paths, as well as source and target data formats.

With a clean and modular code structure, developers can easily extend it to support new file types or additional validation and transformation logic.

**Key Highlights:**

* Written in Go, focusing on performance and portability.
* Validates file existence, size, and JSON structure integrity.
* Easily extensible to handle multiple data formats.
* Provides a simple command-line interface with intuitive flags.

**Example usage:**

```
Terminal:

go build
./cli-convert --input input.json --output output.csv --from json --to csv
```
