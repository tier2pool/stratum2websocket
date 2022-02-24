FROM golang:1.17-alpine AS build

WORKDIR /root/tier2pool

COPY . .

RUN apk add --no-cache gcc musl-dev

RUN mkdir -p ./build && go build -ldflags "-w -s" -o ./build/tier2pool ./command/main.go

FROM alpine:3.15

WORKDIR /root/tier2pool

COPY --from=build /root/tier2pool/build/tier2pool /root/tier2pool/build/tier2pool

ENTRYPOINT ["/root/tier2pool/tier2pool"]
