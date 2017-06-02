VERSION := 0.1

all: build

image:
	docker build -t openpolicyagent/distrib_s3:$(VERSION) .

build: image
