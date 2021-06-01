SHELL := /bin/bash

run:
	go run app/sales-api/main.go --version

tidy:
	go mod tidy
	go mod vendor
