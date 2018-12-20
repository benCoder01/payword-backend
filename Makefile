main: main.go
	go build main.go
	./main

db:
	docker run --name payword-backend -p 27017:27017 -d mongo
