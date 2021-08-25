VERSION?="v1"
root:

build-image:
	docker buildx build --platform linux/amd64 \
			--tag registry.cn-hangzhou.aliyuncs.com/k8s-aa/redis_manger:$(VERSION) \
			-f ./build/Dockerfile \
			.

upload-image:
	docker push registry.cn-hangzhou.aliyuncs.com/k8s-aa/redis_manger:$(VERSION)

release: build-image upload-image
	kubectl delete -f ./deploy/deploy.yaml
	kubectl apply -f ./deploy/deploy.yaml

docker-login:
	docker login --username=$(DOKCER_USERNAME) registry.cn-hangzhou.aliyuncs.com --password $(DOCKER_PASSWORD)