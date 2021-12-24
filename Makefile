.PHONY: build
build:
	go build -o ./build/tier2pool ./main/main.go

build_linux_386:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool_linux_386 ./main/main.go

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool_linux_amd64 ./main/main.go

build_linux_arm:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool_linux_arm ./main/main.go

build_linux_arm64:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool_linux_arm64 ./main/main.go

build_windows_386:
	GOOS=windows GOARCH=amd64 go build -o ./build/tier2pool_windows_386.exe ./main/main.go

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o ./build/tier2pool_windows_amd64.exe ./main/main.go

clear:
	rm -rf ./build
