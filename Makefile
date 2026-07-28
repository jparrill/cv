TEMPLATE ?= tokyo-night
OUTPUT   ?= output

.PHONY: build pdf push clean list-templates

build:
	go run main.go --template $(TEMPLATE) --output $(OUTPUT)

pdf:
	go run main.go --template $(TEMPLATE) --output $(OUTPUT) --pdf

push: build
	git add -A
	git commit -s -m "Update CV"
	git push origin main

clean:
	rm -rf $(OUTPUT)/

list-templates:
	@go run main.go --template _nonexistent_ 2>&1 | grep "Available" || true
