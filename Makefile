all:
	go build -o bin/main cmd/main.go
	PORT=:80 ./bin/main

dev:
	air