build:
	go build -o oreno-fiber main.go

dev:
	go run main.go

run:
	make build && ./oreno-fiber

load:
	cd loadTest && docker compose run k6 run /scripts/loadTest.js && cd -