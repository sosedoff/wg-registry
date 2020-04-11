.PHONY: assets
assets:
	go-assets-builder static \
		-p assets \
		-o assets/assets.go

docker-build:
	docker build -t wg-registry .

linux-binary: docker-build
	docker run -it -d --name=tmp wg-registry bash
	docker cp tmp:/artifacts/wg-registry ./dist/
	docker rm -f tmp