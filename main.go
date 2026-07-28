package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed templates
var templatesFS embed.FS

type Data struct {
	Profile        Profile        `yaml:"profile"`
	Experience     []Experience   `yaml:"experience"`
	Education      []Education    `yaml:"education"`
	Certifications []Certification `yaml:"certifications"`
	Skills         []Skill        `yaml:"skills"`
	Publications   []Publication  `yaml:"publications"`
	Languages      []Language     `yaml:"languages"`
	Volunteering   []Volunteer    `yaml:"volunteering"`
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

func (e Experience) DescriptionHTML() template.HTML {
	return markdownishToHTML(e.Description)
}

type Education struct {
	Institution string `yaml:"institution"`
	Start       string `yaml:"start"`
	End         string `yaml:"end"`
	Degree      string `yaml:"degree"`
}

type Certification struct {
	Name   string `yaml:"name"`
	Issuer string `yaml:"issuer"`
	Date   string `yaml:"date"`
}

type Skill struct {
	Category string   `yaml:"category"`
	Items    []string `yaml:"items"`
}

func (s Skill) ItemsJoined() string {
	return strings.Join(s.Items, ", ")
}

type Publication struct {
	Title     string `yaml:"title"`
	Publisher string `yaml:"publisher"`
	Date      string `yaml:"date"`
	URL       string `yaml:"url"`
}

type Language struct {
	Name  string `yaml:"name"`
	Level string `yaml:"level"`
}

type Volunteer struct {
	Role         string `yaml:"role"`
	Organization string `yaml:"organization"`
	Start        string `yaml:"start"`
	End          string `yaml:"end"`
	Description  string `yaml:"description"`
}

func markdownishToHTML(s string) template.HTML {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	var out strings.Builder
	lines := strings.Split(s, "\n")
	inList := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if inList {
				out.WriteString("</ul>")
				inList = false
			}
			continue
		}

		isBullet := strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ")
		if isBullet {
			if !inList {
				out.WriteString("<ul>")
				inList = true
			}
			out.WriteString("<li>")
			out.WriteString(strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* "))
			out.WriteString("</li>")
		} else {
			if inList {
				out.WriteString("</ul>")
				inList = false
			}
			out.WriteString("<p>")
			out.WriteString(trimmed)
			out.WriteString("</p>")
		}
	}

	if inList {
		out.WriteString("</ul>")
	}

	return template.HTML(out.String())
}

func main() {
	templateName := flag.String("template", "ats-clean", "template name (directory under templates/)")
	dataFile := flag.String("data", "data.yaml", "path to YAML data file")
	outputDir := flag.String("output", "output", "output directory")
	flag.Parse()

	if env := os.Getenv("TEMPLATE"); env != "" && *templateName == "ats-clean" {
		*templateName = env
	}

	raw, err := os.ReadFile(*dataFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", *dataFile, err)
		os.Exit(1)
	}

	var data Data
	if err := yaml.Unmarshal(raw, &data); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	tplPath := fmt.Sprintf("templates/%s/template.html", *templateName)
	tplBytes, err := templatesFS.ReadFile(tplPath)
	if err != nil {
		available := listTemplates()
		fmt.Fprintf(os.Stderr, "Error: template %q not found.\nAvailable templates: %s\n", *templateName, strings.Join(available, ", "))
		os.Exit(1)
	}

	tpl, err := template.New("cv").Parse(string(tplBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering: %v\n", err)
		os.Exit(1)
	}

	outPath := filepath.Join(*outputDir, "cv.html")
	if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Rendered %s template to %s\n", *templateName, outPath)
}

func listTemplates() []string {
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names
}
