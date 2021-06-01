SHELL := /bin/bash

# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"

run:
	go run app/sales-api/main.go

tidy:
	go mod tidy
	go mod vendor
