.PHONY: build

build-consumer:
	docker build -f ./build/Dockerfile.consumer -t arnoldedev/mary-dating-consumer  .
build-seeder:
	docker build -f ./build/Dockerfile.seeder -t arnoldedev/mary-dating-seeder  .

build: build-consumer build-seeder

deploy:
# TODO: some magic with kubectl