.PHONY: build html pdf clean

build:
	go build -o cv main.go

html:
	go run main.go

pdf: html
	# Requires: weasyprint (pip install weasyprint) or wkhtmltopdf
	# Option 1: WeasyPrint (recommended, better CSS support)
	weasyprint output/cv.html output/cv.pdf
	# Option 2: wkhtmltopdf (alternative)
	# wkhtmltopdf output/cv.html output/cv.pdf

clean:
	rm -rf output/
