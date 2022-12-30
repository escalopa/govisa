test:
	go test -v -cover ./...

CFLAGS=-g
export CFLAGS
deploy:
	echo "Deploying image to Docker Hub"
	docker build -f deployments/Dockerfile-telegram -t dekuyo/govisa-telegram:${TAG} --target=production --no-cache .
	docker image push dekuyo/govisa-telegram:${TAG}
.PHONY: test