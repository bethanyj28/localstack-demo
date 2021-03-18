build:
	docker build -t bethanyj28/localstack-demo .
run:
	docker run --rm -p 8080:8080 bethanyj28/localstack-demo
up:
	docker-compose up --build
down:
	docker-compose down
vendor:
	go mod vendor && go mod tidy
