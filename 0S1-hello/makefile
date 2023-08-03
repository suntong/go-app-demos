export GO111MODULE := on

build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build -o ./main

run: build
	./main
