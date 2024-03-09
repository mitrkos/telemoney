run_gsheet:
	go run ./cmd/gsheet_test/main.go

run_tg:
	go run ./cmd/telegram_test/main.go

run_telemoney:
	go run ./cmd/telemoney/main.go

apply_env:
	export $(cat ./.env | xargs) && env


docker_build:
	docker build --progress=plain --tag telemoney .

docker_run_local:
	docker run --env-file .env telemoney

docker_run:
	docker run -e TELEMONEY_TG_BOT_TOKEN -e TELEMONEY_GAUTH_TOKEN telemoney
