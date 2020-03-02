version = $(shell git describe --abbrev=0 --tag)
service = prom-exporter
docker_repo = lukelau/$(service)
.PHONY: docker
docker:
	docker build -t $(docker_repo):$(version) ./ -f Dockerfile
	docker tag $(docker_repo):$(version) $(docker_repo):latest

.PHONY: docker-publish
docker-publish:
	docker push $(docker_repo):$(version)
	docker push $(docker_repo):latest

