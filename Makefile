sqlc:
	docker run --rm -v "$(shell cd)":/src -w /src sqlc/sqlc generate

goose-up:
	backend/bin/goose -dir sql/schema postgres postgres://postgres:OnlyADevPasswOrD@localhost:5432/dbname up

goose-down:
	backend/bin/goose -dir sql/schema postgres postgres://postgres:OnlyADevPasswOrD@localhost:5432/dbname down

air:
	cd backend && air

server:
	cd backend && go build -o cmd/main/main.exe cmd/main/main.go && cmd\main\main.exe

tidy:
	cd backend && go mod tidy

build:
	cd backend && go build -o cmd/main/main.exe cmd/main/main.go