postgres_run:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=azsx0123456 -d postgres:16-alpine
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres16 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:azsx0123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:azsx0123456@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY: createdb, dropdb, migrateup, migratedown, postgres_run