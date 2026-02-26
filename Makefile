MPD_ADDR ?= 127.0.0.1:6600
PORT     ?= 8080

# Run the Go backend + Vite dev server together.
# Use Vite on http://localhost:5173 for HMR during development.
.PHONY: dev
dev:
	@echo "Starting Go backend (MPD at $(MPD_ADDR), API on :$(PORT))"
	@echo "Starting Vite dev server (http://localhost:5173)"
	@echo ""
	$(MAKE) -j2 dev-backend dev-frontend

.PHONY: dev-backend
dev-backend:
	MPD_ADDR=$(MPD_ADDR) PORT=$(PORT) go run .

.PHONY: dev-frontend
dev-frontend:
	cd frontend && npm run dev

# Build everything: frontend + Go binary.
.PHONY: build
build: build-frontend
	go build -v

.PHONY: build-frontend
build-frontend:
	cd frontend && npm run build

# Build for Raspberry Pi (ARMv7).
.PHONY: build-pi
build-pi: build-frontend
	GOOS=linux GOARM=7 GOARCH=arm go build -o mothership

# Run tests.
.PHONY: test
test:
	go test -v -race ./...

# Install frontend dependencies.
.PHONY: install
install:
	cd frontend && npm install

# Clean build artifacts.
.PHONY: clean
clean:
	rm -rf frontend/dist mothership

.PHONY: deploy
deploy:
	./bin/deploy
