.PHONY: network postgres createdb dropdb migrateup migratedown addmigration

network:
	docker network create gn-network

postgres:
	docker run --name postgres12 --network gn-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root gn_db

dropdb:
	docker exec -it postgres12 dropdb gn_db

migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/gn_db?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/gn_db?sslmode=disable" -verbose down

addmigration:
	migrate create -ext sql -dir migrations ${file}