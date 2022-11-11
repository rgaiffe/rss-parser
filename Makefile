NAME="rss-parser"

.PHONY: all dev build test clean

all: build run

dev:
	go run main.go

build:
	go build -o $(NAME) main.go

run: build
	./$(NAME)

test:
	go test -v ./...

clean:
	rm -f $(NAME)

flush-redis:
	docker exec redis-rss redis-cli flushall