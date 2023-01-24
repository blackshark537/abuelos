build:
	go build -o bin/dataprod

run:
	go run ./src/main.go s -p 3000

attack:
	echo "GET http://:3000/api/v1/projection/table/reproductoras/2023" | vegeta attack -rate=10/s -duration=5s | vegeta encode > results.json