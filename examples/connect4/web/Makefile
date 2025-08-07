

# Compile TypeScript entry points to JavaScript using Webpack
build:
	@echo "Compiling TypeScript with Webpack..."
	@# Install webpack dependencies if not present
	@pnpm install
	@# Build with webpack in production mode
	@npx webpack --mode=production
	@echo "✅ Webpack compilation complete!"

# Development build with source maps and watch mode
dev:
	@echo "Starting Webpack development build..."
	@pnpm install
	@npx webpack --mode=development --watch
	
# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf static/gen/js/
	@rm -rf templates/gen/
	@echo "✅ Clean complete!"

.PHONY: build dev clean
