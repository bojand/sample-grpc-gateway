DIST_DIR=./dist

.PHONY: all
all: build

.PHONY: build
build: apigw samplesvc

.PHONY: apigw
apigw:
	cp -R static ./dist
	go build --ldflags="-s -w" -o $(DIST_DIR)/apigw ./cmd/apigw/...

.PHONY: samplesvc
samplesvc:
	go build --ldflags="-s -w" -o $(DIST_DIR)/samplesvc ./cmd/samplesvc/...

.PHONY: clean
clean:
	rm -rf $(DIST_DIR)