.PHONY: binary
.PHONY: vendor

NAME=imgx
PLATFORM=local

binary:
	docker buildx build . -o type=local,dest=bin --platform $(PLATFORM)

install: binary
	cp bin/docker-$(NAME) ~/.docker/cli-plugins/docker-$(NAME)
	chmod +x ~/.docker/cli-plugins/docker-$(NAME)

vendor:
	docker run --init -it --rm \
		-v docker-plugin-$(NAME)-cache:/root/.cache \
		-v ./:/proj -w /proj \
		$(shell docker buildx build . -q --target base) \
			go mod tidy
