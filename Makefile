build:
		docker build -t forum .
run-img:
		docker run --name=forum -p 8080:8080 --rm -d forum
run:
		go run cmd/main.go
stop:
		docker stop forum