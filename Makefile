GITCOMMIT:=$(shell git describe --dirty --always)
GATEWAY:="https://openfaas.weavedx.com"

.PHONY: build
build:
	faas build -f fninfo.yml
	rm -rf build

.PHONY: push
push:
	docker tag stefanprodan/fninfo:latest stefanprodan/fninfo:$(GITCOMMIT)
	docker push stefanprodan/fninfo:$(GITCOMMIT)
	docker push stefanprodan/fninfo:latest

.PHONY: deploy
deploy:
	faas-cli deploy --image=stefanprodan/fninfo:$(GITCOMMIT) --name=fninfo -g=$(GATEWAY)

