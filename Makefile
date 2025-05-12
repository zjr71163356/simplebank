DB_URL=postgresql://root:azsx0123456@localhost:5432/simple_bank?sslmode=disable
image:
	docker build -t simplebank:latest .
run:
	 docker run --network bank-network --name simplebank -p 1234:1234 -e GIN_MODE=release  -e DB_SOURCE=${DB_URL}  simplebank:latest
start:
	docker start simplebank
stop:	
	docker stop simplebank
remove:
	docker stop simplebank
	docker rm simplebank
inspect:
	docker container inspect simplebank
postgres_start:
	docker start postgres16
postgres_run:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=azsx0123456 -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres16 dropdb simple_bank
migrateup:
	migrate -path db/migration -database ${DB_URL} -verbose up
migratedown:
	migrate -path db/migration -database ${DB_URL} -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen --package mockdb -destination db/mock/store.go  github.com/zjr71163356/simplebank/db/sqlc Store
composeup:
	docker compose up
composedown:
	docker compose down	
proto :
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative     --go-grpc_out=pb  --go-grpc_opt=paths=source_relative     proto/*.proto
.PHONY: createdb dropdb migrateup  migratedown  postgres_run server mock image start stop remove inspect