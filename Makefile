.PHONY: build
build:
	go build -o ./build/tier2pool ./main/main.go

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o ./build/tier2pool.exe ./main/main.go

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool ./main/main.go

clear:
	rm -rf ./build
