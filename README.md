# restval-compare

`restval-compare` fetches data from two REST API endpoints, extracts values using JSONPath expressions, and compares them. This tool is useful for verifying data consistency between different API services that might have different response structures.

## Features

- YAML-based configuration
- JSONPath for flexible data extraction
- Custom HTTP headers support
- Verbose logging mode with detailed comparison report

## Installation

```bash
$ go build
```

## Usage

```bash
restval-compare [-verbose] CONFIG_FILE
```

### Arguments

- `CONFIG_FILE`: Path to a YAML configuration file
- `--verbose`: Enable verbose logging (optional)

### Configuration File

Configuration is specified in YAML format. See `config.yaml.example` for a template.

## About JSONPath

`restval-compare` uses the [theory/jsonpath](https://github.com/theory/jsonpath) library, which implements the RFC 9535 JSONPath standard.

For more details, refer to the [RFC 9535 JSONPath specification](https://www.rfc-editor.org/rfc/rfc9535.html).

## Exit Status

- 0: All values match
- 1: Mismatch found or error occurred

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
