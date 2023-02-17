build:
	go build -o bin/abuelos

run:
	go run ./src/main.go s -p 3000

attack:
	echo "GET http://:3000/api/v1/projection/table/reproductoras/2023" | vegeta attack -rate=10/s -duration=5s | vegeta encode > results.json

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/macos/abuelos

dk-build:
	docker build -t abuelos:alpine .
dk-run:
	docker run --name dtp_abuelos -e DB_NAME=abuelos -e DB_URL=mongodb://mongodb:27017 -p 3000:3000 -d --net dtprod abuelos

dk-run-mongo:
	docker run --name mongodb -e ALLOW_EMPTY_PASSWORD=yes -p 27017 -d --net dtprod bitnami/mongodb