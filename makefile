SHELL := /bin/bash

# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"

# ==============================================================================
# Building containers

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

# ==============================================================================

run:
	go run app/sales-api/main.go

tidy:
	go mod tidy
	go mod vendor
