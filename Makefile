TEMPLATE ?= ats-clean
OUTPUT   ?= output

.PHONY: build html pdf clean list-templates

build:
	go build -o cv main.go

html:
	go run main.go --template $(TEMPLATE) --output $(OUTPUT)

pdf: html
	@echo "PDF generation not yet implemented (chromedp coming soon)"

clean:
	rm -rf $(OUTPUT)/

list-templates:
	@go run main.go --template _nonexistent_ 2>&1 | grep "Available" || true
