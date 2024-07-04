run:
	@go run cmd/main.go

test:
	@go test ./...

test-coverage:
	@go test -v -coverprofile cover.out ./...
	@go tool cover -html=cover.out -o cover.html
	@open cover.html

generate-mock:
	@echo "GENERATING ..."
	@echo "[API]"
	@echo "DONE ..."

migrate-up-local:
	@echo "Executing local database migration ..."
	@goose -dir migrations mysql "root:@tcp(127.0.0.1:3306)/dpo-test?parseTime=true" up