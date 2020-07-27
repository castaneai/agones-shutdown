IMAGE_BASE=gameserver

build:
	make build-v1
	make build-v2

build-v1:
	@eval $$(minikube docker-env) ;\
	docker build --build-arg VERSION=v1 -t ${IMAGE_BASE}:v1 .

build-v2:
	@eval $$(minikube docker-env) ;\
	docker build --build-arg VERSION=v2 -t ${IMAGE_BASE}:v2 .
