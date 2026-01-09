.PHONY: dev build generate clean

# Development with hot reload
dev:
	@echo "Starting development server with hot reload..."
	air

# Build for production
build: generate
	@echo "Building production binary..."
	go build -o musiq-server .

# Generate templ and tailwind
generate:
	@echo "Generating templ files..."
	templ generate
	@echo "Building Tailwind CSS..."
	~/go/bin/tailwindcss -i web/static/css/input.css -o web/static/css/output.css --minify

# Generate without minify (for development)
generate-dev:
	@echo "Generating templ files..."
	templ generate
	@echo "Building Tailwind CSS..."
	~/go/bin/tailwindcss -i web/static/css/input.css -o web/static/css/output.css

# Clean build artifacts
clean:
	rm -rf tmp/
	rm -f musiq-server
	rm -f web/static/css/output.css

# Run the server directly
run: generate-dev
	go run .
