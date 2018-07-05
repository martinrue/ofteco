build: assets
	@go build -o ./dist/frekvenco ./cmd/frekvenco

build-linux: assets
	@GOOS=linux GOARCH=amd64 go build -o ./dist/frekvenco-linux-amd64 ./cmd/frekvenco

assets:
	@esc -ignore=".go|.sketch" -prefix=assets -o ./assets/assets.go -pkg assets ./assets

clean:
	@rm -rf ./dist

tools:
	@go install ./vendor/github.com/mjibson/esc

.PHONY: build build-linux assets clean tools
