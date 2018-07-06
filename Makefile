build: assets
	@go build -o ./dist/ofteco ./cmd/ofteco

build-linux: assets
	@GOOS=linux GOARCH=amd64 go build -o ./dist/ofteco-linux-amd64 ./cmd/ofteco

assets:
	@esc -ignore=".go|.sketch" -prefix=assets -o ./assets/assets.go -pkg assets ./assets

clean:
	@rm -rf ./dist

tools:
	@go install ./vendor/github.com/mjibson/esc

.PHONY: build build-linux assets clean tools
