BACKEND_BINARY=backendApp
GAME_BINARY=gameApp


## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

up_build: build_backend build_game
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build_backend: builds the backend binary as a linux executable
build_backend:
	@echo "Building backend binary..."
	cd ../backend-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BACKEND_BINARY} .
	@echo "Done!"

## build_game: builds the game binary as a linux executable
build_game:
	@echo "Building backend binary..."
	cd ../game-service && env GOOS=linux CGO_ENABLED=0 go build -o ${GAME_BINARY} .
	@echo "Done!"
