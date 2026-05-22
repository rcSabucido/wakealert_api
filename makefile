-include .env
export

up:
	docker compose up -d 

down:
	docker compose down 

reset:
	docker compose down -v && docker compose up -d 

migrate-up:
	migrate -path internal/db/migrations -database "${DATABASE_URL}" up

migrate-down:
	migrate -path internal/db/migrations -database "${DATABASE_URL}" down 1

migrate-new:
	migrate create -ext sql -dir internal/db/migrations -seq ${name}

sqlc:
	sqlc generate

run:
	go run ./cmd/api

psql: 
	psql "${DATABASE_URL}"