stop:
	docker compose stop
build: stop
	docker compose build && docker compose up -d
