build:
	docker build -t bethanyj28/localstack-demo .
run:
	docker run --rm -p 8080:8080 bethanyj28/localstack-demo
vendor:
	go mod vendor && go mod tidy
