e2e.up:
	docker compose --file docker-compose.e2e.yaml up \
					--remove-orphans \
					--build \
					--detach

e2e.down:
	docker compose --file docker-compose.e2e.yaml down

e2e.test:
	gotest -tags=e2e ./... -v

e2e: proto build e2e.up migrate.up e2e.test migrate.down e2e.down
