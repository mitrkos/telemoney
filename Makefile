run_gsheet:
	go run ./cmd/gsheet_test/main.go

run_tg:
	go run ./cmd/telegram_test/main.go

run_telemoney:
	go run ./cmd/telemoney/main.go

apply_env:
	export $(cat ./.env | xargs) && env

docker-run:
	docker compose up -d

docker-run-dev:
	docker compose up --build

docker-down:
	docker compose down

