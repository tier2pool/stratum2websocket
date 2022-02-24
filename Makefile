.PHONY: build
build:
	go build -o ./build/tier2pool ./command/main.go

build_linux_386:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool_linux_386 ./command/main.go

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool_linux_amd64 ./command/main.go

build_linux_arm:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool_linux_arm ./command/main.go

build_linux_arm64:
	GOOS=linux GOARCH=amd64 go build -o ./build/tier2pool_linux_arm64 ./command/main.go

build_windows_386:
	GOOS=windows GOARCH=amd64 go build -o ./build/tier2pool_windows_386.exe ./command/main.go

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o ./build/tier2pool_windows_amd64.exe ./command/main.go

build_image:
	docker build -t tier2pool/tier2pool:v0.1.0 -t tier2pool/tier2pool:latest .

clear:
	rm -rf ./build
