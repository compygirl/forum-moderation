build:
		docker build -t forum .
run-img:
		docker run --name=forum -p 8082:8082 --rm -d forum
run:
		go run cmd/main.go
stop:
		docker stop forum