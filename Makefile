TOKEN?=none
PORT?=8080
CHAT_ID?=none

BUILD_ARGS=--build-arg TOKEN=$(TOKEN) --build-arg PORT=$(PORT) --build-arg CHAT_ID=$(CHAT_ID)

docker_build: 
	@if [ $(TOKEN) == "none" ] ; then \
		echo "TOKEN not set!"; \
		exit 1;\
	fi
	@if [ $(CHAT_ID) == "none" ] ; then \
		echo "CHAT_ID not set!"; \
		exit 1;\
	fi
	
	docker build  ${BUILD_ARGS} -t s1kai/memes_bot:0.0.1 .
