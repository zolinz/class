SHELL := /bin/bash

# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"

# // To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem

# make admin
# eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQyNzQ1NTQuNzM4OTI4LCJpYXQiOjE2MjI3Mzg1NTQuNzM4OTMyLCJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiIxMjM0NTY3Nzg5IiwiUm9sZXMiOlsiQURNSU4iXX0.btlXINXk0co2munSzWNQtumxzkAz4dKvsjTNHLzQpa7lsbCUvIDo8gQZ9Z5r6QXJAHihuhJfAX66SOibcQIymE-UD_Iu_Hy0HbYpkYm3WZ0NgHHQAdSjqCqqxmWUby6ERErrYNwttf38HDk-MaCoIE5LCmutPSIPFmtBlQLT3aR-EAdKOv9Odpt1j8JzHdGkY8V42HlH55SjGMRo_-e-kcwWYtiojfU9vbRI5uf0Z0tdgLaoSEuKUedVvmq-9P8_tog9Hu4FJMJGrWiLPIRFAbhTVDmxz_Zkn8SKv3CxAiqrL7H_WHu9jM2urvfvCb5tItrwWMOxPRFYQOocGlkUKA
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/auth
# curl http://localhost:3000/auth

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
# Running from within k8s/dev

kind-up:
	kind create cluster --image kindest/node:v1.20.2 --name ardan-starter-cluster --config zarf/k8s/dev/kind-config.yaml

kind-down:
	kind delete cluster --name ardan-starter-cluster

kind-load:
	kind load docker-image sales-api-amd64:1.0 --name ardan-starter-cluster

kind-services:
	kustomize build zarf/k8s/dev | kubectl apply -f -

kind-logs:
	kubectl logs -lapp=sales-api --all-containers=true -f --tail=100

kind-status:
	kubectl get nodes
	kubectl get pods --watch

kind-update: sales-api
	kind load docker-image sales-api-amd64:1.0 --name ardan-starter-cluster
	kubectl delete pods -lapp=sales-api

# ==============================================================================

run:
	go run app/sales-api/main.go

admin:
	go run app/admin/main.go

tidy:
	go mod tidy
	go mod vendor
