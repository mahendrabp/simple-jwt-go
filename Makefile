hello:
	echo "Hello"

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/simple-jwt-go?sslmode=disable" -verbose up

swagger-generate:
	@swag init -g ./cmd/main.go
