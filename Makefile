.PHONY: all build clean dev

all: build

# Build frontend, copy dist into backend/public/, compile Go binary
build:
	cd frontend && npm install && npm run build
	rm -rf backend/public
	cp -r frontend/dist backend/public
	cd backend && go build -ldflags="-s -w" -o pkl-server .

# Remove build artifacts
clean:
	rm -rf backend/public backend/pkl-server

# Run dev servers (frontend on :5173, backend on :8082)
dev:
	@echo "Start backend:  cd backend && go run ."
	@echo "Start frontend: cd frontend && npm run dev"
