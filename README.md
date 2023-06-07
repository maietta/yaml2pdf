# yaml2pdf
A simple, cross-platform command line utility to generate PDF files from HTML templates and YAML files containing data.

A use case for this tool is to generate invoices from YAML files containing invoice data and HTML templates containing the invoice layout. Can be used for any other use case where you need to generate PDF files from HTML templates and data files.

This program is designed to be used in a processing pipeline, for example to generate invoices from a web application. It is not designed to be used as a standalone application itself-- although it can be used that way.

## Build the binary

```bash
go build -o yaml2pdf.exe cmd/yaml2pdf/main.go
```

## How to run after building the binary
```bash
yaml2pdf.exe --data data/inv-001.yaml --template template.html --output pdfs/
```

# How to run in development mode
```bash
go run cmd/yaml2pdf/main.go --data data/inv-001.yaml --template template.html --output pdfs/
```

Like the project? Give it a star :star: and spread the word!

Looking for something more custom? [Hire me](https://maietta.consulting) :wink: