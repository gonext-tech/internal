build:
	tailwindcss -i view/css/styles.css -o public/styles.css
	@templ generate view
	@go build -o bin/blog main.go

test:
	@go test -v ./...

run: build
	@./bin/blog

tailwind:
	./tailwindcss -i views/css/styles.css -o public/styles.css --watch

css:
	tailwindcss -i views/css/styles.css -o public/styles.css --watch

templ:
	@templ generate -watch -proxy=http://localhost:7001
