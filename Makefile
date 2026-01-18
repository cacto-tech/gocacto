.PHONY: help install dev build run templ clean artisan css css-watch

GOPATH_BIN := $(shell go env GOPATH)/bin

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest
	@if ! command -v node >/dev/null 2>&1; then \
		echo "⚠️  Node.js not found. Installing Node.js dependencies skipped."; \
		echo "   Install Node.js from https://nodejs.org/ to use Tailwind CSS v4"; \
	else \
		echo "📦 Installing Node.js dependencies for Tailwind CSS v4..."; \
		npm install; \
		echo "✅ Node.js dependencies installed"; \
	fi
	@echo "✅ Dependencies installed"

templ: ## Generate templ files
	@echo "🎨 Generating templ components..."
	$(GOPATH_BIN)/templ generate
	@echo "✅ Templ generation complete"

css: ## Build Tailwind CSS v4 (production)
	@echo "🎨 Building Tailwind CSS v4..."
	@if command -v npm >/dev/null 2>&1; then \
		npm run css:build; \
		echo "✅ CSS built (only used classes included)"; \
	else \
		echo "⚠️  npm not found. Skipping CSS build."; \
		echo "   Install Node.js to build Tailwind CSS"; \
	fi

css-watch: ## Watch Tailwind CSS v4 (development)
	@echo "👀 Watching Tailwind CSS v4..."
	@if command -v npm >/dev/null 2>&1; then \
		npm run css:watch; \
	else \
		echo "⚠️  npm not found. Install Node.js to use Tailwind CSS watch mode."; \
	fi

dev: templ ## Run with hot reload
	@echo "🔥 Starting development server with hot reload..."
	@echo "💡 Tip: Run 'make css-watch' in another terminal to watch CSS changes"
	$(GOPATH_BIN)/air

build: templ css ## Build the application (with CSS)
	@echo "🔨 Building..."
	go build -o cacto-cms ./cmd/server
	go build -o artisan ./cmd/artisan
	@echo "✅ Build complete"

run: templ ## Run the application
	@echo "🚀 Starting server..."
	go run ./cmd/server/main.go

artisan: ## Run artisan CLI (usage: make artisan ARGS="migrate:fresh --seed")
	@echo "🛠️  Running artisan..."
	go run ./cmd/artisan/main.go $(ARGS)

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -f cacto-cms artisan
	rm -rf tmp/
	find . -name "*_templ.go" -delete
	rm -f *.db *.db-shm *.db-wal
	rm -f web/static/css/output.css
	@echo "✅ Cleaned"

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

tidy: ## Tidy go modules
	@echo "📦 Tidying modules..."
	go mod tidy

.DEFAULT_GOAL := help
