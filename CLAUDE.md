# CV Renderer

YAML-driven CV renderer with multi-template support and PDF generation via headless Chrome (chromedp). Single Go binary, no external dependencies beyond Chrome for PDF.

## Architecture

```
data.yaml -> Go renderer -> embedded HTML templates -> output/cv.html + output/cv.pdf
```

- `data.yaml`: single source of truth for all CV content (profile, experience, education, certs, skills, publications, languages, volunteering)
- `main.go`: reads YAML, selects template, renders HTML via `html/template`, optionally generates PDF with chromedp
- `templates/`: each subdirectory contains a `template.html` that receives the full `Data` struct
- Templates are embedded at compile time via `//go:embed templates`

## Templates

| Name | Description |
|------|-------------|
| `tokyo-night` | **Default and recommended.** Dark theme inspired by Tokyo Night color scheme |
| `ats-clean` | ATS-friendly clean layout, light theme |
| `minimal` | Basic layout (needs updating to current struct -- uses old field names) |

### Adding a new template

1. Create `templates/<name>/template.html`
2. Template receives the `Data` struct (see type definitions in `main.go`)
3. Use `{{.Profile.Name}}`, `{{range .Experience}}`, etc.
4. Experience descriptions support markdown-ish syntax (bullets with `- ` or `* `, paragraphs) via `.DescriptionHTML` method
5. Skills expose `.ItemsJoined()` for comma-separated output
6. Build with `make build TEMPLATE=<name>` to test

## Quick Actions

- `make build` -- render HTML + PDF (always generates both)
- `make push` -- build + commit + push to GitHub, triggers deploy to Pages

## Important: Always Generate PDF

The Tokyo Night template has a "Download PDF" button linking to `cv.pdf`. If the PDF is stale, users download an outdated version. `make build` includes `--pdf` by default -- never skip PDF generation.

## Build

```bash
make build                        # HTML + PDF with tokyo-night (default)
make build TEMPLATE=ats-clean     # Use a different template
make push                         # Build + git add + commit -s + push
make clean                        # Remove output/
make list-templates               # Show available templates
```

### Flags (direct `go run`)

```bash
go run main.go --template tokyo-night --data data.yaml --output output --pdf
```

- `--template`: template directory name (default: `ats-clean`, Makefile overrides to `tokyo-night`)
- `--data`: path to YAML data file (default: `data.yaml`)
- `--output`: output directory (default: `output`)
- `--pdf`: generate PDF in addition to HTML (requires Chrome/Chromium)
- `TEMPLATE` env var overrides `--template` when the flag is at its default

## Deploy

GitHub Actions (`.github/workflows/deploy.yml`) on push to `main`:
1. Builds with `tokyo-night` template into `_site/`
2. Copies `cv.html` to `index.html`
3. Deploys to GitHub Pages at `jparrill.github.io/cv`

## Conventions

- Sign commits: `git commit -s`
- Go tests: table-driven with Gherkin naming (`When X it should Y`)
- Commit messages in English
- This repo is part of the `jparrill.github.io` ecosystem
