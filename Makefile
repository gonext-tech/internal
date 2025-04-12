build:
	@templ generate view
	@go build -o bin/internal cmd/main.go
	@pm2 restart pm2.config.js

test:
	@go test -v ./...

run: build
	@./bin/internal

tailwind:
	./tailwindcss -i views/css/styles.css -o public/css/styles.css --watch

templ:
	@templ generate -watch -proxy=http://localhost:7001
