.PHONY: e2e
e2e: e2e.up e2e.test e2e.down

.PHONY: e2e.up
e2e.up:
	docker compose --file compose.backend.e2e.yaml up \
					--remove-orphans \
					--build \
					--detach

.PHONY: e2e.down
e2e.down:
	docker compose --file compose.backend.e2e.yaml down

.PHONY: e2e.test
e2e.test:
	gotest -tags=e2e ./... -v
