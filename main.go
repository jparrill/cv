package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Data is the top-level YAML structure.
type Data struct {
	Profile      Profile      `yaml:"profile"`
	Experience   []Experience `yaml:"experience"`
	Education    []Education  `yaml:"education"`
	Certifications []Certification `yaml:"certifications"`
	Skills       []Skill      `yaml:"skills"`
	Languages    []Language   `yaml:"languages"`
	Volunteering []Volunteer  `yaml:"volunteering"`
}

type Profile struct {
	Name     string `yaml:"name"`
	Headline string `yaml:"headline"`
	Location string `yaml:"location"`
	Email    string `yaml:"email"`
	Github   string `yaml:"github"`
	Linkedin string `yaml:"linkedin"`
	About    string `yaml:"about"`
}

type Experience struct {
	Company     string `yaml:"company"`
	Location    string `yaml:"location"`
	Title       string `yaml:"title"`
	Start       string `yaml:"start"`
	End         string `yaml:"end"`
	Current     bool   `yaml:"current"`
	Description string `yaml:"description"`
}

type Education struct {
	Institution string `yaml:"institution"`
	Start       string `yaml:"start"`
	End         string `yaml:"end"`
	Degree      string `yaml:"degree"`
}

type Certification struct {
	Name string `yaml:"name"`
	Date string `yaml:"date"`
}

type Skill struct {
	Category string   `yaml:"category"`
	Items    []string `yaml:"items"`
}

type Language struct {
	Name  string `yaml:"name"`
	Level string `yaml:"level"`
}

type Volunteer struct {
	Role        string `yaml:"role"`
	Organization string `yaml:"organization"`
	Start       string `yaml:"start"`
	End         string `yaml:"end"`
	Description string `yaml:"description"`
}

func main() {
	dataFile := "data.yaml"
	templateFile := "templates/cv.html"
	outputDir := "output"

	// Parse YAML
	raw, err := os.ReadFile(dataFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", dataFile, err)
		os.Exit(1)
	}

	var data Data
	if err := yaml.Unmarshal(raw, &data); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	// Read template
	tplBytes, err := os.ReadFile(templateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading template: %v\n", err)
		os.Exit(1)
	}

	// Parse and execute template
	tpl, err := template.New("cv").Parse(string(tplBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// Write HTML
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering template: %v\n", err)
		os.Exit(1)
	}

	outPath := filepath.Join(outputDir, "cv.html")
	if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Rendered to %s\n", outPath)
}
