build:
	go build -o bin/dataprod

run:
	go run ./src/main.go s -p 3000

attack:
	echo "GET http://:3000/api/v1/projection/table/reproductoras/2023" | vegeta attack -rate=10/s -duration=5s | vegeta encode > results.json

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/macos/dataprod

dk-build:
	docker build -t dataprod:latest .
dk-run:
	docker run --name dtprod -e DB_NAME=dataprod -e DB_URL=mongodb://localhost:27017 -p 3000:3000 dataprod