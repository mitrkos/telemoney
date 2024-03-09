run_gsheet:
	go run ./cmd/gsheet_test/main.go

run_tg:
	go run ./cmd/telegram_test/main.go

run_telemoney:
	go run ./cmd/telemoney/main.go

run_tests:
	go test -v ./...

run_test_build:
	go build -v -o ./telemoney ./cmd/telemoney
	rm ./telemoney

apply_env:
	export $(cat ./.env | xargs) && env

docker-run:
	docker compose up -d

docker-run-dev:
	docker compose up --build

docker-run-cicd:
	docker compose -f docker-compose-test.yml up --build  

docker-down:
	docker compose down

