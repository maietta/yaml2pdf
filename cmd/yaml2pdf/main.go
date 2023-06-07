package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/playwright-community/playwright-go"
)

func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func printUsage() {
	fmt.Println("Generate a PDF from a template and YAML data file.")
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {

	// Get the --data flag.
	dataPtr := flag.String("data", "data.yaml", "a filename ending with .yaml")
	// Get the --template flag.
	templatePtr := flag.String("template", "template.html", "the template filename")
	// Get the --output flag.
	outputPtr := flag.String("output", ".", "the output directory")
	flag.Usage = printUsage
	flag.Parse()

	// Check if any flags are passed.
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	// Open the data file.
	dataFile, err := os.Open(*dataPtr)
	assertErrorToNilf("could not open data file: %w", err)
	defer func() {
		if err := dataFile.Close(); err != nil {
			log.Fatalf("could not close data file: %v", err)
		}
	}()

	// Read the YAML data.
	yamlData, err := ioutil.ReadAll(dataFile)
	assertErrorToNilf("could not read data file: %w", err)

	// Unmarshal the YAML data into a map.
	var yamlStruct map[string]interface{}
	err = yaml.Unmarshal([]byte(yamlData), &yamlStruct)
	if err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return
	}

	// Open the template file.
	templateFile, err := os.Open(*templatePtr)
	assertErrorToNilf("could not open template file: %w", err)
	defer func() {
		if err := templateFile.Close(); err != nil {
			log.Fatalf("could not close template file: %v", err)
		}
	}()

	// Read the template file content.
	templateContent, err := ioutil.ReadAll(templateFile)
	assertErrorToNilf("could not read template file: %w", err)

	// Parse the template content.
	t, err := template.New("template").Parse(string(templateContent))
	assertErrorToNilf("could not parse template: %w", err)

	// Convert the *dataPtr to a regular string
	dataPath := *dataPtr

	// Remove the .yaml extension from the dataPath
	dataPath = dataPath[:len(dataPath)-5]

	// Create the output file path.
	outputPath := filepath.Join(*outputPtr, fmt.Sprintf("%s.html", filepath.Base(dataPath)))
	outputFile, err := os.Create(outputPath)
	assertErrorToNilf("could not create output file: %w", err)
	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Fatalf("could not close output file: %v", err)
		}
	}()

	// Execute the template with the artwork data and write to the output file.
	err = t.Execute(outputFile, yamlStruct)
	assertErrorToNilf("could not execute template: %w", err)

	pw, err := playwright.Run()
	assertErrorToNilf("could not launch playwright: %w", err)
	defer func() {
		if err := pw.Stop(); err != nil {
			log.Fatalf("could not stop Playwright: %v", err)
		}
	}()

	// Launch browser.
	browser, err := pw.Chromium.Launch()
	assertErrorToNilf("could not launch Chromium: %w", err)
	defer func() {
		if err := browser.Close(); err != nil {
			log.Fatalf("could not close browser: %v", err)
		}
	}()

	// Create context.
	context, err := browser.NewContext()
	assertErrorToNilf("could not create context: %w", err)
	defer func() {
		if err := context.Close(); err != nil {
			log.Fatalf("could not close context: %v", err)
		}
	}()

	// Create page.
	page, err := context.NewPage()
	assertErrorToNilf("could not create page: %w", err)
	defer func() {
		if err := page.Close(); err != nil {
			log.Fatalf("could not close page: %v", err)
		}
	}()

	// Get current directory.
	dir, err := os.Getwd()
	assertErrorToNilf("could not get current directory: %w", err)

	// Convert to browser-compatible file URL.
	fileURL := fmt.Sprintf("file://%s/%s", dir, outputPath)

	// Navigate to file URL.
	if _, err = page.Goto(fileURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	// Create PDF.
	pdfPath := filepath.Join(*outputPtr, fmt.Sprintf("%s.pdf", filepath.Base(dataPath)))
	_, err = page.PDF(playwright.PagePdfOptions{
		Path:  playwright.String(pdfPath),
		Scale: playwright.Float(0.9), // Set to 0.9 to scale to 90%.
	})
	assertErrorToNilf("could not create PDF: %w", err)

	fmt.Printf("Template executed successfully. Output file: %s\n", outputPath)
}
