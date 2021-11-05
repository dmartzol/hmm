
up:
	go mod tidy
	go mod download
	docker-compose up --remove-orphans -d --build

down:
	docker compose -p hmm down
