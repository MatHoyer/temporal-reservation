.PHONY: web-build
web-build:
	cd web && npm install && npm run build
	rm -rf internal/webassets/dist
	mkdir -p internal/webassets/dist
	cp -r web/dist/. internal/webassets/dist/

.PHONY: run-api
run-api:
	go run ./cmd/api

.PHONY: run-worker
run-worker:
	go run ./cmd/worker

.PHONY: run-web
run-web:
	cd web && npm run dev

.PHONY: test
test:
	go test ./... -race -cover

.PHONY: lint
lint:
	go vet ./...

.PHONY: compose-up
compose-up:
	docker compose -f docker/docker-compose.yml up --build

.PHONY: compose-down
compose-down:
	docker compose -f docker/docker-compose.yml down -v
