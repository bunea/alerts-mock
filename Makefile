run:
	go run main.go

docker-build:
	docker build --tag 'alerts-mock' .

docker-run:
	docker run -p 1323:1323 'alerts-mock'

dev: docker-build docker-run
